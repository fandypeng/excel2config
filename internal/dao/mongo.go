package dao

import (
	"context"
	"encoding/json"
	pb "excel2config/api"
	"excel2config/internal/def"
	"excel2config/internal/helper"
	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strconv"
	"time"
)

const (
	dbname                = "sheets"
	tableExcelList        = "excel_list"
	tableUserList         = "user_list"
	tableTokenList        = "token_list"
	tableGroupList        = "group_list"
	tableExportRecordList = "export_record_list"
)

func NewMongo() (m *mongo.Client, cf func(), err error) {
	var (
		ct  paladin.Map
		cfg struct {
			Dsn string
		}
	)
	if err = paladin.Get("mongo.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Dsn))
	if err != nil {
		return nil, func() {}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, func() {}, err
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Errorw(ctx, "err", err, "msg", "mongo db ping failed")
		return
	}
	cf = func() { client.Disconnect(context.Background()) }
	m = client
	return
}

func (d *dao) PingMongo(ctx context.Context) error {
	if err := d.mongo.Ping(ctx, readpref.Primary()); err != nil {
		log.Errorw(ctx, "err", err, "msg", "mongo db ping failed")
		return err
	}
	return nil
}

func (d *dao) LoadExcel(ctx context.Context, gridKey string) (sheets []*model.Sheet, err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	//先查出status=1的sheet的完整内容，其他status=0的sheet只取出基础信息，不包含sheet内容
	singleRes := c.FindOne(ctx, bson.D{{"status", 1}})
	activeSheet := new(model.Sheet)
	err = singleRes.Decode(activeSheet)
	if err != nil {
		log.Errorw(ctx, "err", err, "msg", "mongo decode error")
		return nil, err
	}
	opt := options.Find()
	opt.SetProjection(bson.D{{"celldata", false}})
	opt.SetSort(bson.D{{"order", 1}})
	corsor, err := c.Find(ctx, bson.D{{"status", 0}, {"deleted", bson.M{"$not": bson.M{"$eq": 1}}}}, opt)
	if err != nil {
		return
	}
	sheets = make([]*model.Sheet, 0)
	err = corsor.All(ctx, &sheets)
	if err != nil {
		log.Errorw(ctx, "err", err, "msg", "mango decode error")
		return
	}
	sheets = append(sheets, activeSheet)
	return
}

func (d *dao) LoadExcelSheet(ctx context.Context, gridKey string, indexs []string) (sheets map[string][]model.Cell, err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filters := make([]bson.M, 0)
	filters = append(filters, bson.M{"deleted": 0})
	for _, index := range indexs {
		filters = append(filters, bson.M{"index": index})
	}
	opt := options.Find()
	opt.SetProjection(bson.D{{"celldata", true}, {"index", true}})
	corsor, err := c.Find(ctx, bson.M{"$or": filters}, opt)
	if err != nil {
		return
	}
	sheetInfos := make([]*model.Sheet, 0)
	err = corsor.All(ctx, &sheetInfos)
	if err != nil {
		log.Errorw(ctx, "err", err, "msg", "mango decode error")
		return
	}
	sheets = make(map[string][]model.Cell)
	for _, sheet := range sheetInfos {
		if sheet.Celldata == nil {
			sheet.Celldata = make([]model.Cell, 0)
		}
		sheets[sheet.Index] = sheet.Celldata
	}
	return
}

func (d *dao) LoadAllSheet(ctx context.Context, gridKey string) (sheets []*model.Sheet, err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	opt := options.Find()
	opt.SetSort(bson.D{{"order", 1}})
	corsor, err := c.Find(ctx, bson.D{{"name", bson.M{"$not": bson.M{"$eq": def.DefaultIntroductionSheet}}}}, opt)
	if err != nil {
		return
	}
	sheets = make([]*model.Sheet, 0)
	err = corsor.All(ctx, &sheets)
	if err != nil {
		log.Errorw(ctx, "err", err, "msg", "mango decode error")
		return
	}
	return
}

func (d *dao) LoadSheetByName(ctx context.Context, gridKey, sheetName string) (sheet *model.Sheet, err error) {
	if sheetName == def.DefaultIntroductionSheet {
		err = ecode.Int(int(def.ErrSheetName))
		return
	}
	c := d.mongo.Database(dbname).Collection(gridKey)
	sheet = new(model.Sheet)
	res := c.FindOne(ctx, bson.D{{"name", bson.M{"$eq": sheetName}}})
	err = res.Decode(sheet)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Errorw(ctx, "err", err, "msg", "mango decode error")
		}
		return
	}
	return
}

func (d *dao) UpdateGridValue(ctx context.Context, gridKey string, req *model.UpdateV) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.I, "celldata": bson.M{"$elemMatch": bson.M{"r": req.R, "c": req.C}}}
	formatV, err := d.format2Bson(req.V)
	if err != nil {
		log.Errorw(ctx, "err", err, "v", req.V, "msg", "format2bson error")
		return
	}
	res, err := c.UpdateOne(ctx, filter, bson.D{{"$set", bson.D{{"celldata.$.v", formatV}}}})
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	if res.ModifiedCount <= 0 {
		formatCell, err := d.format2Bson(req.Cell)
		if err != nil {
			log.Errorw(ctx, "err", err, "cell", req.Cell, "msg", "format2bson error")
			return err
		}
		_, err = c.UpdateOne(ctx, bson.M{"index": req.I}, bson.M{"$push": bson.M{"celldata": formatCell}})
	}
	return err
}

func (d *dao) UpdateGridMulti(ctx context.Context, gridKey string, req *model.UpdateRV) (err error) {
	subIndex := 0
	cs, ce := req.Range.Column[0], req.Range.Column[1]
	rs, re := req.Range.Row[0], req.Range.Row[1]
	for c := cs; c <= ce; c++ {
		index := 0
		for r := rs; r <= re; r++ {
			customReq := &model.UpdateV{
				Cell: model.Cell{
					C: model.FlexInt(c),
					R: model.FlexInt(r),
					V: req.V[index][subIndex],
				},
				I: req.I,
			}
			log.Errorw(ctx, "c", c+1, "r", r+1, "v", req.V[index][subIndex], "msg", "update grid")
			err = d.UpdateGridValue(ctx, gridKey, customReq)
			if err != nil {
				log.Errorw(ctx, "err", err, "req", customReq, "msg", "update error")
				return
			}
			index++
		}
		subIndex++
	}
	return err
}

func (d *dao) UpdateGridConfig(ctx context.Context, gridKey string, req *model.UpdateCG) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.I}
	formatV, err := d.format2Bson(map[string]interface{}{
		"config." + req.K: req.V,
	})
	if err != nil {
		log.Errorw(ctx, "err", err, "v", req.V, "msg", "format2bson error")
		return
	}
	_, err = c.UpdateOne(ctx, filter, bson.D{{"$set", formatV}})
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) UpdateGridCommon(ctx context.Context, gridKey string, req *model.UpdateCommon) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.I}
	formatV, err := d.format2Bson(map[string]interface{}{
		req.K: req.V,
	})
	if err != nil {
		log.Errorw(ctx, "err", err, "v", req.V, "msg", "format2bson error")
		return
	}
	_, err = c.UpdateOne(ctx, filter, bson.D{{"$set", formatV}})
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) UpdateCalcChain(ctx context.Context, gridKey string, req *model.UpdateCalcChain) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.I}
	cc := new(model.CalcChain)
	err = json.Unmarshal([]byte(req.V), cc)
	if err != nil {
		log.Errorw(ctx, "err", err, "v", req.V, "msg", "json unmarshal error")
		return
	}
	formatV, err := d.format2Bson(cc)
	if err != nil {
		log.Errorw(ctx, "err", err, "v", req.V, "msg", "format2bson error")
		return
	}
	switch req.Op {
	case "add":
		_, err = c.UpdateOne(ctx, filter, bson.D{{"$push", bson.M{"calcChain": formatV}}})
	case "update":
		_, err = c.UpdateOne(ctx, filter, bson.D{{"$set", bson.M{"calcChain." + strconv.Itoa(req.Pos): formatV}}})
	case "del":
		_, err = c.UpdateOne(ctx, filter, bson.D{{"$unset", bson.M{"calcChain." + strconv.Itoa(req.Pos): 1}}})
		if err != nil {
			_, err = c.UpdateOne(ctx, filter, bson.D{{"$pull", bson.M{"calcChain": nil}}})
		}
	}
	_, err = c.UpdateOne(ctx, filter, bson.D{{"$set", formatV}})
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) UpdateRowColumn(ctx context.Context, gridKey string, req *model.UpdateRowColumn) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.I}
	switch req.T {
	case "drc": //删除行或列
		drs := req.V.Index
		dre := req.V.Index + req.V.Len
		switch req.RC {
		case "r": //行操作
			_, err = c.UpdateOne(ctx, filter, bson.M{"$pull": bson.M{"celldata": bson.M{"r": bson.M{"$gte": drs, "$lt": dre}}}})
			if err == nil { //剩余行，修改行号
				opt := &options.UpdateOptions{ArrayFilters: &options.ArrayFilters{Filters: make([]interface{}, 0)}}
				opt.ArrayFilters.Filters = append(opt.ArrayFilters.Filters, bson.M{"elem.r": bson.M{"$gte": dre}})
				_, err = c.UpdateMany(ctx, filter, bson.M{"$inc": bson.M{"celldata.$[elem].r": -req.V.Len}}, opt)
			}
			if err == nil { //总行号修改
				_, err = c.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"row": -req.V.Len}})
			}
		case "c": //列操作
			_, err = c.UpdateOne(ctx, filter, bson.M{"$pull": bson.M{"celldata": bson.M{"c": bson.M{"$gte": drs, "$lt": dre}}}})
			if err == nil { //剩余列，修改列数
				opt := &options.UpdateOptions{ArrayFilters: &options.ArrayFilters{Filters: make([]interface{}, 0)}}
				opt.ArrayFilters.Filters = append(opt.ArrayFilters.Filters, bson.M{"elem.c": bson.M{"$gte": dre}})
				_, err = c.UpdateMany(ctx, filter, bson.M{"$inc": bson.M{"celldata.$[elem].c": -req.V.Len}}, opt)
			}
			if err == nil { //总列数修改
				_, err = c.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"column": -req.V.Len}})
			}
		}
	case "arc": //增加行或列
		drs := req.V.Index
		switch req.RC {
		case "r": //行操作
			_, err = c.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"row": req.V.Len}})
			if err == nil { //剩余行，修改行号
				opt := &options.UpdateOptions{ArrayFilters: &options.ArrayFilters{Filters: make([]interface{}, 0)}}
				opt.ArrayFilters.Filters = append(opt.ArrayFilters.Filters, bson.M{"elem.r": bson.M{"$gte": drs}})
				_, err = c.UpdateMany(ctx, filter, bson.M{"$inc": bson.M{"celldata.$[elem].r": req.V.Len}}, opt)
			}
		case "c": //列操作
			_, err = c.UpdateOne(ctx, filter, bson.M{"$inc": bson.M{"column": req.V.Len}})
			if err == nil { //剩余列，修改列数
				opt := &options.UpdateOptions{ArrayFilters: &options.ArrayFilters{Filters: make([]interface{}, 0)}}
				opt.ArrayFilters.Filters = append(opt.ArrayFilters.Filters, bson.M{"elem.c": bson.M{"$gte": drs}})
				_, err = c.UpdateMany(ctx, filter, bson.M{"$inc": bson.M{"celldata.$[elem].c": req.V.Len}}, opt)
			}
		}
	}
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) UpdateFilter(ctx context.Context, gridKey string, req *model.UpdateFilter) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.I}
	if req.V == nil {
		_, err = c.UpdateOne(ctx, filter, bson.D{{"$unset", bson.M{"filter": "", "filter_select": ""}}})
	} else {
		formatFilter, _ := d.format2Bson(req.V.Filter)
		formatFilterSelect, _ := d.format2Bson(req.V.FilterSelect)
		_, err = c.UpdateOne(ctx, filter, bson.D{{"$set", bson.M{"filter": formatFilter, "filter_select": formatFilterSelect}}})
	}
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) AddSheet(ctx context.Context, gridKey string, req *model.AddSheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	formatBson, _ := d.format2Bson(req.V)
	_, err = c.InsertOne(ctx, formatBson)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) UpdateSheet(ctx context.Context, gridKey string, sheet *model.Sheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	//先判断sheet是否存在
	res := c.FindOne(ctx, bson.M{"name": sheet.Name})
	if res.Err() == mongo.ErrNoDocuments {
		sheet.Index = "sheet_" + strconv.Itoa(int(time.Now().Unix()))
		err = d.AddSheet(ctx, gridKey, &model.AddSheet{
			V: sheet,
		})
	} else {
		oldSheet := &model.Sheet{}
		err = res.Decode(&oldSheet)
		if err != nil {
			return
		}
		sheet.Status = oldSheet.Status
		sheet.Order = oldSheet.Order
		formatBson, _ := d.format2Bson(sheet)
		_, err = c.UpdateMany(ctx, bson.M{"name": sheet.Name}, bson.M{"$set": formatBson})
	}
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "sheet", sheet, "msg", "update error")
		return
	}
	return err
}

func (d *dao) CopySheet(ctx context.Context, gridKey string, req *model.CopySheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	//先查出原来sheet的数据，修改index和name之后，添加到db中
	singleRes := c.FindOne(ctx, bson.D{{"index", req.V.CopyIndex}})
	activeSheet := new(model.Sheet)
	err = singleRes.Decode(activeSheet)
	if err != nil {
		log.Errorw(ctx, "err", err, "msg", "mongo decode error")
		return err
	}
	activeSheet.Name = req.V.Name
	activeSheet.Index = req.I
	formatBson, _ := d.format2Bson(activeSheet)
	_, err = c.InsertOne(ctx, formatBson)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) DeleteSheet(ctx context.Context, gridKey string, req *model.DeleteSheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.V.DeleteIndex}
	_, err = c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"deleted": 1}})
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) ExcelInfo(ctx context.Context, gridKey string) (excelInfo *model.Excel, err error) {
	c := d.mongo.Database(dbname).Collection(tableExcelList)
	_id, err := primitive.ObjectIDFromHex(gridKey)
	if err != nil {
		return
	}
	filter := bson.M{"_id": _id}
	singleRes := c.FindOne(ctx, filter)
	excelInfo = new(model.Excel)
	err = singleRes.Decode(excelInfo)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "msg", "mongo decode error")
		return
	}
	return
}

func (d *dao) ExcelInfoByGroupId(ctx context.Context, gid, name string) (excelInfo *model.Excel, err error) {
	c := d.mongo.Database(dbname).Collection(tableExcelList)
	filter := bson.M{"group_id": gid, "name": name}
	singleRes := c.FindOne(ctx, filter)
	excelInfo = new(model.Excel)
	err = singleRes.Decode(excelInfo)
	if err != nil {
		log.Errorw(ctx, "err", err, "gid", gid, "msg", "mongo decode error")
		return
	}
	return
}

func (d *dao) RecoverSheet(ctx context.Context, gridKey string, req *model.RecoverSheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"index": req.V.RecoverIndex}
	_, err = c.UpdateOne(ctx, filter, bson.M{"$unset": bson.M{"deleted": ""}})
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) UpdateSheetOrder(ctx context.Context, gridKey string, req *model.UpdateSheetOrder) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	for index, order := range req.V {
		filter := bson.M{"index": index}
		_, err = c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"order": order}})
	}
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) ToggleSheet(ctx context.Context, gridKey string, req *model.ToggleSheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"status": 1}
	_, err = c.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"status": 0}})
	if err == nil {
		filter := bson.M{"index": req.V}
		_, err = c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"status": 1}})
	}
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) HideOrShowSheet(ctx context.Context, gridKey string, req *model.HideOrShowSheet) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	switch req.Op {
	case "hide":
		filter := bson.M{"index": req.I}
		_, err = c.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"hide": req.V, "status": 0}})
		if err == nil {
			filter := bson.M{"index": req.Cur}
			_, err = c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"status": 1}})
		}
	case "show":
		filter := bson.M{"status": 1}
		_, err = c.UpdateMany(ctx, filter, bson.M{"$set": bson.M{"status": 0}})
		if err == nil {
			filter := bson.M{"index": req.I}
			_, err = c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"status": 1, "hide": req.V}})
		}
	}

	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "req", req, "msg", "update error")
		return
	}
	return err
}

func (d *dao) RemDeletedSheet(ctx context.Context, gridKey string) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.M{"deleted": 1}
	_, err = c.DeleteMany(ctx, filter)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "msg", "delete error")
		return
	}
	return err
}

func (d *dao) CreateExcel(ctx context.Context, uid, gridKey, remark, groupId string) (eid string, err error) {
	lc := d.mongo.Database(dbname).Collection(tableExcelList)
	excel := model.Excel{
		Name:         gridKey,
		CreateTime:   time.Now().Unix(),
		EditTime:     time.Now().Unix(),
		Owner:        uid,
		Remark:       remark,
		GroupId:      groupId,
		Contributers: []string{uid},
	}
	formatBson, _ := d.format2Bson(excel)
	res, err := lc.InsertOne(ctx, formatBson)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "msg", "update error")
		return
	}
	excelId := res.InsertedID.(primitive.ObjectID).Hex()
	err = d.mongo.Database(dbname).CreateCollection(ctx, excelId)
	if err != nil {
		return
	}
	c := d.mongo.Database(dbname).Collection(excelId)
	sheet := &model.Sheet{
		Celldata: helper.GetIntroduceCellData(),
		Column:   50,
		Index:    "sheet_" + strconv.Itoa(int(time.Now().Unix())),
		Name:     def.DefaultIntroductionSheet,
		Order:    1,
		Row:      100,
		Status:   1,
		Config:   helper.GetIntroduceSheetConfig(),
	}
	formatBson, _ = d.format2Bson(sheet)
	_, err = c.InsertOne(ctx, formatBson)
	if err != nil {
		log.Errorw(ctx, "err", err, "gridKey", gridKey, "sheet", sheet, "msg", "update error")
		return
	}
	eid = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (d *dao) UpdateExcel(ctx context.Context, id, remark string, contributers []string) (err error) {
	c := d.mongo.Database(dbname).Collection(tableExcelList)
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	update := bson.D{}
	if len(contributers) > 0 {
		ctrs, err := d.format2Bson(contributers)
		if err == nil {
			update = append(update, bson.E{Key: "contributers", Value: ctrs})
		}
	}
	if len(remark) > 0 {
		update = append(update, bson.E{Key: "remark", Value: remark})
	}
	_, err = c.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": update})
	if err != nil {
		return
	}
	return
}

func (d *dao) DeleteExcel(ctx context.Context, gridKey string) (err error) {
	c := d.mongo.Database(dbname).Collection(tableExcelList)
	_id, err := primitive.ObjectIDFromHex(gridKey)
	if err != nil {
		return
	}
	_, err = c.DeleteOne(ctx, bson.M{"_id": _id})
	if err == nil {
		err = d.mongo.Database(dbname).Collection(gridKey).Drop(ctx)
	}
	if err != nil {
		return
	}
	return
}

func (d *dao) ExcelList(ctx context.Context, lastTime int64, limit int64, groupId string) (list []*pb.SimpleExcel, err error) {
	c := d.mongo.Database(dbname).Collection(tableExcelList)
	filter := bson.M{"create_time": bson.M{"$gte": lastTime}, "group_id": groupId}
	opt := &options.FindOptions{
		Limit: &limit,
		Sort:  bson.M{"create_time": 1},
	}
	cur, err := c.Find(ctx, filter, opt)
	if err != nil {
		log.Errorw(ctx, "err", err, "lastTime", lastTime, "msg", "get excel list error")
		return
	}
	list = make([]*pb.SimpleExcel, 0)
	err = cur.All(ctx, &list)
	if err != nil {
		log.Errorw(ctx, "err", err, "lastTime", lastTime, "msg", "decode excel list error")
		return
	}
	return
}

func (d *dao) format2Bson(v interface{}) (doc bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}
