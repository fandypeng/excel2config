package dao

import (
	"context"
	"excel2config/internal/model"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const dbname = "sheets"

func NewMongo() (m *mongo.Client, cf func(), err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://fandy:fandypeng@127.0.0.1:27017/?authSource=sheets"))
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
		log.With("err", err).Error("mongo db ping failed")
		return
	}
	cf = func() { client.Disconnect(context.Background()) }
	m = client
	return
}

func (d *dao) PingMongo(ctx context.Context) error {
	if err := d.mongo.Ping(ctx, readpref.Primary()); err != nil {
		log.With("err", err).Error("mongo db ping failed")
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
		log.With("err", err).Errorln("mongo decode error")
		return nil, err
	}
	opt := options.Find()
	opt.SetProjection(bson.D{{"name", true}, {"index", true}, {"order", true}, {"status", true}})
	opt.SetSort(bson.D{{"order", 1}})
	corsor, err := c.Find(ctx, bson.D{{"status", 0}}, opt)
	if err != nil {
		return
	}
	sheets = make([]*model.Sheet, 0)
	err = corsor.All(ctx, &sheets)
	if err != nil {
		log.With("err", err).Errorln("mango decode error")
		return
	}
	sheets = append(sheets, activeSheet)
	return
}

func (d *dao) LoadExcelSheet(ctx context.Context, gridKey string, indexs []string) (sheets map[string][]model.Cell, err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filters := make([]bson.M, 0)
	for _, index := range indexs {
		filters = append(filters, bson.M{"index": index})
	}
	opt := options.Find()
	opt.SetProjection(bson.D{{"celldata", true}, {"name", true}})
	corsor, err := c.Find(ctx, bson.M{"$or": filters}, opt)
	if err != nil {
		return
	}
	sheetInfos := make([]*model.Sheet, 0)
	err = corsor.All(ctx, &sheetInfos)
	if err != nil {
		log.With("err", err).Errorln("mango decode error")
		return
	}
	sheets = make(map[string][]model.Cell)
	for _, sheet := range sheetInfos {
		sheets[sheet.Name] = sheet.Celldata
	}
	return
}

func (d *dao) UpdateGridValue(ctx context.Context, gridKey string, req *model.UpdateGridReq) (err error) {
	c := d.mongo.Database(dbname).Collection(gridKey)
	filter := bson.D{{"r", req.R}, {"c", req.C}, {"index", req.I}}
	_, err = c.UpdateOne(ctx, filter, bson.D{{"$set", bson.D{{"v", req.V}}}})
	if err != nil {
		log.With("err", err).With("gridKey", gridKey).With("req", req).Errorln("update error")
	}
	return err
}
