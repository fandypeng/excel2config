package dao

import (
	"context"
	"excel2config/internal/model"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *dao) GroupList(ctx context.Context, uid, gid string) (groupList []*model.GroupInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableGroupList)
	filter := bson.M{"members": bson.M{"$elemMatch": bson.M{"uid": uid}}, "IsDev": true}
	if len(gid) > 0 {
		groupId, err := primitive.ObjectIDFromHex(gid)
		if err == nil {
			filter = bson.M{"members": bson.M{"$elemMatch": bson.M{"uid": uid}}, "_id": groupId}
		}
	}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return
	}
	groupList = make([]*model.GroupInfo, 0)
	err = cursor.All(ctx, &groupList)
	if err != nil {
		log.With("err", err).Errorln("mango decode error")
		return
	}
	return
}

func (d *dao) GroupAdd(ctx context.Context, groupInfo *model.GroupInfo) (groupId string, err error) {
	c := d.mongo.Database(dbname).Collection(tableGroupList)
	format, err := d.format2Bson(groupInfo)
	if err != nil {
		log.With("err", err).With("groupInfo", groupInfo).Errorln("format group error")
		return
	}
	res, err := c.InsertOne(ctx, format)
	if err != nil {
		log.With("err", err).With("groupInfo", groupInfo).Errorln("insert group error")
		return
	}
	groupId = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (d *dao) GroupUpdate(ctx context.Context, groupInfo *model.GroupInfo) (err error) {
	c := d.mongo.Database(dbname).Collection(tableGroupList)
	gid := groupInfo.Gid
	groupInfo.Gid = ""
	format, err := d.format2Bson(groupInfo)
	if err != nil {
		log.With("err", err).With("groupInfo", groupInfo).Errorln("format group error")
		return
	}
	groupId, err := primitive.ObjectIDFromHex(gid)
	if err != nil {
		return
	}
	_, err = c.UpdateOne(ctx, bson.M{"_id": groupId}, bson.M{"$set": format})
	if err != nil {
		log.With("err", err).With("groupInfo", groupInfo).Errorln("update group error")
		return
	}
	return
}

func (d *dao) GroupInfo(ctx context.Context, groupId string) (groupInfo *model.GroupInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableGroupList)
	gid, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		return
	}
	groupInfo = &model.GroupInfo{}
	res := c.FindOne(ctx, bson.M{"_id": gid})
	err = res.Decode(groupInfo)
	if err != nil {
		log.With("err", err).With("groupId", groupId).Errorln("get group error")
		return
	}
	return
}

func (d *dao) AddExportRecord(ctx context.Context, record *model.ExportRecord) (err error) {
	c := d.mongo.Database(dbname).Collection(tableExportRecordList)
	format, err := d.format2Bson(record)
	if err != nil {
		return
	}
	_, err = c.InsertOne(ctx, format)
	if err != nil {
		return
	}
	//每张表最多保留50条历史记录
	//c.DeleteMany(ctx)
	return
}

func (d *dao) GetExportRecordList(ctx context.Context, gridKey, sheetName string) (list []model.ExportRecord, err error) {
	c := d.mongo.Database(dbname).Collection(tableExportRecordList)
	filter := bson.M{"gridKey": gridKey, "sheetName": sheetName}
	limit := int64(30)
	opt := &options.FindOptions{
		Limit: &limit,
		Sort:  bson.M{"time": -1},
	}
	cursor, err := c.Find(ctx, filter, opt)
	if err != nil {
		return
	}
	list = make([]model.ExportRecord, 0)
	err = cursor.All(ctx, &list)
	if err != nil {
		log.With("err", err).Errorln("mango decode error")
		return
	}
	return
}

func (d *dao) GetExportRecord(ctx context.Context, gridKey, sheetName, recordId string) (record *model.ExportRecord, err error) {
	c := d.mongo.Database(dbname).Collection(tableExportRecordList)
	rid, err := primitive.ObjectIDFromHex(recordId)
	if err != nil {
		log.With("err", err).Errorln("mango decode error")
		return
	}
	filter := bson.M{"gridKey": gridKey, "sheetName": sheetName, "_id": rid}
	res := c.FindOne(ctx, filter)
	record = &model.ExportRecord{}
	err = res.Decode(&record)
	if err != nil {
		log.With("err", err).Errorln("mango decode error")
		return
	}
	return
}
