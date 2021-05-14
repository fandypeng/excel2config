package dao

import (
	"context"
	"errors"
	"excel2config/internal/def"
	"excel2config/internal/helper"
	"excel2config/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (d *dao) GetUser(ctx context.Context, email string) (userInfo *model.UserInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableUserList)
	res := c.FindOne(ctx, bson.M{"email": email})
	userInfo = new(model.UserInfo)
	err = res.Decode(userInfo)
	return
}

func (d *dao) GetUserInfos(ctx context.Context, uids []string) (userInfos map[string]*model.UserInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableUserList)
	userList := make([]*model.UserInfo, 0)
	objIds := make([]primitive.ObjectID, 0)
	for _, uid := range uids {
		oid, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			continue
		}
		objIds = append(objIds, oid)
	}
	cur, err := c.Find(ctx, bson.M{"_id": bson.M{"$in": objIds}})
	if err != nil {
		return
	}
	err = cur.All(ctx, &userList)
	if err != nil {
		return
	}
	userInfos = make(map[string]*model.UserInfo)
	for _, userInfo := range userList {
		userInfos[userInfo.Uid] = userInfo
	}
	return
}

func (d *dao) GetUserByUid(ctx context.Context, uid string) (userInfo *model.UserInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableUserList)
	oid, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return
	}
	res := c.FindOne(ctx, bson.M{"_id": oid})
	userInfo = new(model.UserInfo)
	err = res.Decode(userInfo)
	return
}

func (d *dao) AddUser(ctx context.Context, email, passwd, name string, loginType int, openId string) (userInfo *model.UserInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableUserList)
	userInfo = &model.UserInfo{
		UserName:  name,
		Email:     email,
		Passwd:    helper.Md5Sum(passwd + def.PasswdSalt),
		RegTime:   time.Now().Unix(),
		Avatar:    helper.GetRandomAvatar(),
		LoginType: loginType,
		OpenId:    openId,
	}
	buser, err := d.format2Bson(userInfo)
	if err != nil {
		return
	}
	inRes, err := c.InsertOne(ctx, buser)
	if err != nil {
		return
	}
	userInfo.Uid = inRes.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (d *dao) SaveUser(ctx context.Context, userInfo *model.UserInfo) (err error) {
	c := d.mongo.Database(dbname).Collection(tableUserList)
	uid, err := primitive.ObjectIDFromHex(userInfo.Uid)
	if err != nil {
		return
	}
	var formatInfo = *userInfo
	formatInfo.Uid = ""
	buser, err := d.format2Bson(formatInfo)
	if err != nil {
		return
	}
	_, err = c.UpdateOne(ctx, bson.M{"_id": uid}, bson.M{"$set": buser})
	if err != nil {
		return
	}
	return
}

func (d *dao) SaveToken(ctx context.Context, uid, token string) (err error) {
	c := d.mongo.Database(dbname).Collection(tableTokenList)
	filter := bson.M{"uid": uid}
	upres, err := c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"token": token, "update_time": time.Now().Unix()}})
	if err != nil {
		return
	}
	if upres.ModifiedCount == 0 {
		_, err = c.InsertOne(ctx, bson.M{"uid": uid, "token": token, "update_time": time.Now().Unix()})
		if err != nil {
			return
		}
	}
	return
}

func (d *dao) GetToken(ctx context.Context, uid string) (token string, err error) {
	c := d.mongo.Database(dbname).Collection(tableTokenList)
	filter := bson.M{"uid": uid}
	res := c.FindOne(ctx, filter)
	var tokenInfo = new(model.Token)
	if err = res.Decode(tokenInfo); err != nil {
		return
	}
	if time.Now().Unix()-tokenInfo.UpdateTime > def.DaySeconds {
		err = errors.New("token expired")
	}
	token = tokenInfo.Token
	return
}

func (d *dao) GetUserInfosByKeyword(ctx context.Context, keyword string) (userInfos map[string]*model.SimpleUserInfo, err error) {
	c := d.mongo.Database(dbname).Collection(tableUserList)
	userList := make([]*model.UserInfo, 0)
	opt := options.Find()
	opt.SetProjection(bson.D{{"_id", true}, {"username", true}, {"avatar", true}})
	cur, err := c.Find(ctx, bson.M{"username": bson.M{"$regex": keyword, "$options": "i"}}, opt)
	if err != nil {
		return
	}
	err = cur.All(ctx, &userList)
	if err != nil {
		return
	}
	userInfos = make(map[string]*model.SimpleUserInfo)
	for _, userInfo := range userList {
		userInfos[userInfo.Uid] = &model.SimpleUserInfo{
			Uid:      userInfo.Uid,
			UserName: userInfo.UserName,
			Avatar:   userInfo.Avatar,
		}
	}
	return
}
