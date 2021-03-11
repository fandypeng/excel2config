package service

import (
	"context"
	"encoding/json"
	pb "excel2config/api"
	"excel2config/internal/def"
	"excel2config/internal/helper"
	"excel2config/internal/model"
	"excel2config/internal/server/sessions"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strconv"
	"time"
)

func (s *Service) GroupList(ctx context.Context, req *pb.GroupListReq) (resp *pb.GroupListResp, err error) {
	sess := sessions.Default(ctx.(*bm.Context))
	var uid string
	uidInterface := sess.Get("uid")
	if uidInterface != nil {
		uid = uidInterface.(string)
	} else {
		err = ecode.Int(int(def.ErrNeedLogin))
		return
	}
	list, err := s.dao.GroupList(ctx, uid)
	if err != nil {
		return
	}
	formatList := make([]*pb.GroupInfo, 0)
	for _, ginfo := range list {
		formatList = append(formatList, s.copyGroupInfo(ctx, ginfo))
	}
	resp = &pb.GroupListResp{
		List: formatList,
	}
	return
}

func (s *Service) GroupAdd(ctx context.Context, req *pb.AddGroupReq) (resp *pb.AddGroupResp, err error) {
	sess := sessions.Default(ctx.(*bm.Context))
	var uid string
	uidInterface := sess.Get("uid")
	if uidInterface != nil {
		uid = uidInterface.(string)
	} else {
		err = ecode.Int(int(def.ErrNeedLogin))
		return
	}
	groupInfo := &model.GroupInfo{
		Name:    req.Name,
		Avatar:  req.Avatar,
		Remark:  req.Remark,
		Store:   []int32{},
		AddTime: time.Now().Unix(),
		Owner:   uid,
		Members: []model.SimpleUserInfo{{Uid: uid}},
	}
	gid, err := s.dao.GroupAdd(ctx, groupInfo)
	if err != nil {
		return
	}
	groupInfo.Gid = gid
	resp = &pb.AddGroupResp{
		GroupInfo: s.copyGroupInfo(ctx, groupInfo),
	}
	return
}
func (s *Service) GroupUpdate(ctx context.Context, req *pb.UpdateGroupReq) (resp *pb.UpdateGroupResp, err error) {
	sess := sessions.Default(ctx.(*bm.Context))
	var uid string
	uidInterface := sess.Get("uid")
	if uidInterface != nil {
		uid = uidInterface.(string)
	} else {
		err = ecode.Int(int(def.ErrNeedLogin))
		return
	}
	groupInfo, err := s.dao.GroupInfo(ctx, req.Id)
	if err != nil {
		return
	}
	if groupInfo.Owner != uid {
		err = ecode.Int(int(def.ErrPermissionDenied))
		return
	}
	members := make([]model.SimpleUserInfo, 0)
	uids := make([]string, 0)
	for _, minfo := range req.GroupInfo.Members {
		if helper.Contains(uids, minfo.Uid) {
			continue
		}
		members = append(members, model.SimpleUserInfo{
			Uid:  minfo.Uid,
			Role: minfo.Role,
		})
		uids = append(uids, minfo.Uid)
	}
	groupInfo = &model.GroupInfo{
		Gid:            groupInfo.Gid,
		Name:           req.GroupInfo.Name,
		Avatar:         req.GroupInfo.Avatar,
		Remark:         req.GroupInfo.Remark,
		Store:          req.GroupInfo.Store,
		AddTime:        req.GroupInfo.AddTime,
		Owner:          req.GroupInfo.Owner,
		Members:        members,
		RedisDSN:       req.GroupInfo.RedisDSN,
		RedisPassword:  req.GroupInfo.RedisPassword,
		RedisKeyPrefix: req.GroupInfo.RedisKeyPrefix,
		MysqlDSN:       req.GroupInfo.MysqlDSN,
		MongodbDSN:     req.GroupInfo.MongodbDSN,
	}
	err = s.dao.GroupUpdate(ctx, groupInfo)
	if err != nil {
		return
	}
	resp = &pb.UpdateGroupResp{}
	return
}

func (s *Service) TestConnection(ctx context.Context, req *pb.TestConnectionReq) (resp *pb.TestConnectionResp, err error) {
	resp = &pb.TestConnectionResp{Connected: 1}
	switch req.DsnType {
	case def.DsnTypeRedis:
		redisConf := helper.ParseRedisDsn(req.Dsn, req.Pwd)
		rc := redis.NewRedis(redisConf)
		defer rc.Close()
		_, connErr := rc.Do(ctx, "SET", "ping", "pong")
		if connErr != nil {
			log.Errorw(ctx, "redis conn error", connErr)
			resp.Connected = 0
			return
		}
	case def.DsnTypeMysql:
		mysqlConf := helper.ParseMysqlDsn(req.Dsn)
		db := sql.NewMySQL(mysqlConf)
		defer db.Close()
		if connErr := db.Ping(ctx); connErr != nil {
			log.Errorw(ctx, "mysql conn error, connErr: ", connErr)
			resp.Connected = 0
			return
		}
	case def.DsnTypeMongodb:
		client, connErr := mongo.NewClient(options.Client().ApplyURI(req.Dsn))
		if connErr != nil {
			resp.Connected = 0
			log.Errorw(ctx, "mongodb new config error, connErr: ", connErr)
			return
		}
		connErr = client.Connect(ctx)
		if connErr != nil {
			resp.Connected = 0
			log.Errorw(ctx, "mongodb conn error, connErr: ", connErr)
			return
		}
		if connErr := client.Ping(ctx, readpref.Primary()); connErr != nil {
			log.Errorw(ctx, "mongodb ping error, connErr: ", connErr)
			client.Disconnect(ctx)
			resp.Connected = 0
			return
		}
		client.Disconnect(ctx)
	}
	return
}

func (s *Service) ExportConfigToDB(ctx context.Context, req *pb.ExportConfigToDBReq) (resp *pb.ExportConfigToDBResp, err error) {
	excelInfo, err := s.dao.ExcelInfo(ctx, req.GridKey)
	if err != nil {
		log.Errorw(ctx, "ws conn gridKey not exist, gridKey: ", req.GridKey)
		err = ecode.Int(int(def.ErrTableNotExist))
		return
	}
	gid := excelInfo.GroupId
	groupInfo, err := s.dao.GroupInfo(ctx, gid)
	if err != nil {
		log.Errorw(context.TODO(), "ws conn group not exist, groupId: ", gid)
		err = ecode.Int(int(def.ErrGroupNotExist))
		return
	}
	if len(groupInfo.Store) == 0 {
		err = ecode.Int(int(def.ErrGroupStoreEmpty))
		return
	}
	sheet, err := s.dao.LoadSheetByName(ctx, req.GridKey, req.SheetName)
	if err != nil {
		log.Errorw(ctx, "get sheet failed, err: ", err)
		return
	}
	formatInfo, err := sheet.Format()
	if err != nil {
		return
	}
	for _, dsnType := range groupInfo.Store {
		switch dsnType {
		case def.DsnTypeRedis:
			redisConf := helper.ParseRedisDsn(groupInfo.RedisDSN, groupInfo.RedisPassword)
			rc := redis.NewRedis(redisConf)
			defer rc.Close()
			key := groupInfo.RedisKeyPrefix + req.SheetName
			_, connErr := rc.Do(ctx, "SET", "ping", "pong")
			if connErr == nil {
				jsonbytes, _ := json.Marshal(formatInfo.Content)
				_, connErr = rc.Do(ctx, "set", key, string(jsonbytes))
			}
			if connErr == nil {
				_, connErr = rc.Do(ctx, "publish", groupInfo.RedisKeyPrefix+def.DefaultRedisPubsubChannel, key)
			}
			if connErr != nil {
				log.Errorw(ctx, "redis conn error", connErr)
				err = ecode.Int(int(def.ErrGroupExportRedisFailed))
				return
			}
		case def.DsnTypeMysql:
			mysqlConf := helper.ParseMysqlDsn(groupInfo.MysqlDSN)
			db := sql.NewMySQL(mysqlConf)
			defer db.Close()
			if connErr := db.Ping(ctx); connErr != nil {
				log.Errorw(ctx, "mysql conn error, connErr: ", connErr)
				err = ecode.Int(int(def.ErrGroupGetConfigFailed))
				return
			}
			// build a temp table to test the format
			testTableName := req.SheetName + "_" + strconv.Itoa(int(time.Now().Unix()))
			err = s.exportTableToMysql(ctx, db, formatInfo, testTableName)
			if err == nil {
				err = s.exportTableToMysql(ctx, db, formatInfo, req.SheetName)
			}
			if err == nil {
				tx, txerr := db.Begin(ctx)
				if txerr != nil {
					log.Errorw(ctx, "mysql begin error, err: ", txerr)
					err = txerr
					return
				}
				s.dropTable(tx, testTableName)
				err = tx.Commit()
			}
			if err != nil {
				log.Errorw(ctx, "mysql conn error", err)
				err = ecode.Int(int(def.ErrGroupExportMysqlFailed))
				return
			}
		case def.DsnTypeMongodb:
			//client, connErr := mongo.NewClient(options.Client().ApplyURI(groupInfo.MongodbDSN))
			//if connErr != nil {
			//	log.Errorw(ctx, "mongodb new config error, connErr: ", connErr)
			//	return
			//}
			//connErr = client.Connect(ctx)
			//if connErr != nil {
			//	log.Errorw(ctx, "mongodb conn error, connErr: ", connErr)
			//	return
			//}
			//if connErr := client.Ping(ctx, readpref.Primary()); connErr != nil {
			//	log.Errorw(ctx, "mongodb ping error, connErr: ", connErr)
			//	client.Disconnect(ctx)
			//	return
			//}
			//filter := bson.M{}
			//collections, err := client.Database(def.DefaultMongoDataBaseName).ListCollectionNames(ctx, filter)
			//if err != nil {
			//	log.Errorw(ctx, "mongodb list collections error, err: ", err)
			//	client.Disconnect(ctx)
			//	return
			//}
			//if !helper.Contains(collections, req.SheetName) {
			//	err = client.Database(def.DefaultMongoDataBaseName).CreateCollection(ctx, req.SheetName)
			//	if err != nil {
			//		log.Errorw(ctx, "mongodb list collections error, err: ", err)
			//		client.Disconnect(ctx)
			//		return
			//	}
			//}
			//c := client.Database(def.DefaultMongoDataBaseName).Collection(req.SheetName)
			//client.Disconnect(ctx)
		}
	}
	return
}

func (s *Service) GetConfigFromDB(ctx context.Context, req *pb.GetConfigFromDBReq) (resp *pb.GetConfigFromDBResp, err error) {
	excelInfo, err := s.dao.ExcelInfo(ctx, req.GridKey)
	if err != nil {
		log.Errorw(ctx, "ws conn gridKey not exist, gridKey: ", req.GridKey)
		err = ecode.Int(int(def.ErrTableNotExist))
		return
	}
	gid := excelInfo.GroupId
	groupInfo, err := s.dao.GroupInfo(ctx, gid)
	if err != nil {
		log.Errorw(context.TODO(), "ws conn group not exist, groupId: ", gid)
		err = ecode.Int(int(def.ErrGroupNotExist))
		return
	}
	if len(groupInfo.Store) == 0 {
		err = ecode.Int(int(def.ErrGroupStoreEmpty))
		return
	}
	resp = &pb.GetConfigFromDBResp{}
	for _, dsnType := range groupInfo.Store {
		switch dsnType {
		case def.DsnTypeRedis:
			redisConf := helper.ParseRedisDsn(groupInfo.RedisDSN, groupInfo.RedisPassword)
			rc := redis.NewRedis(redisConf)
			defer rc.Close()
			key := groupInfo.RedisKeyPrefix + req.SheetName
			_, connErr := rc.Do(ctx, "SET", "ping", "pong")
			var jsonstr string
			if connErr == nil {
				jsonstr, connErr = redis.String(rc.Do(ctx, "get", key))
			}
			if connErr != nil && connErr != redis.ErrNil {
				log.Errorw(ctx, "redis conn error", connErr)
				err = ecode.Int(int(def.ErrGroupGetConfigFailed))
				return
			}
			resp.Jsonstr = jsonstr
			return
		case def.DsnTypeMysql:
			mysqlConf := helper.ParseMysqlDsn(groupInfo.MysqlDSN)
			db := sql.NewMySQL(mysqlConf)
			defer db.Close()
			if connErr := db.Ping(ctx); connErr != nil {
				log.Errorw(ctx, "mysql conn error, connErr: ", connErr)
				err = ecode.Int(int(def.ErrGroupGetConfigFailed))
				return
			}
			res := make([]interface{}, 0)
			rows, connErr := db.Query(ctx, "select * from "+req.SheetName)
			if connErr == nil {
				cols, _ := rows.Columns()
				for rows.Next() {
					var row = make([]interface{}, len(cols))
					var rowp = make([]interface{}, len(cols))
					for i, _ := range row {
						rowp[i] = &row[i]
					}
					connErr = rows.Scan(rowp...)
					if connErr != nil {
						break
					}
					data := make(map[string]interface{})
					for i := 0; i < len(cols); i++ {
						columnName := cols[i]
						columnValue := *rowp[i].(*interface{})
						strval := string(columnValue.([]byte))
						data[columnName] = strval
						if intval, err := strconv.Atoi(strval); err == nil {
							data[columnName] = intval
						}
					}
					res = append(res, data)
				}
			}
			if connErr == nil {
				connErr = rows.Err()
			}
			if connErr != nil && connErr != sql.ErrNoRows {
				log.Errorw(ctx, "mysql conn error", connErr)
				return
			}
			jsonbytes, _ := json.Marshal(res)
			resp.Jsonstr = string(jsonbytes)
			return
		}
	}
	return
}

func (s *Service) copyGroupInfo(ctx context.Context, groupInfo *model.GroupInfo) *pb.GroupInfo {
	members := make([]*pb.SimpleUserInfo, 0)
	uids := make([]string, 0)
	for _, minfo := range groupInfo.Members {
		uids = append(uids, minfo.Uid)
	}
	userInfos, err := s.dao.GetUserInfos(ctx, uids)
	if err != nil {
		log.Errorw(ctx, "get user infos error", " err: ", err, " uids: ", uids)
	}
	for _, minfo := range groupInfo.Members {
		userInfo, ok := userInfos[minfo.Uid]
		if !ok {
			continue
		}
		members = append(members, &pb.SimpleUserInfo{
			Uid:      minfo.Uid,
			UserName: userInfo.UserName,
			Avatar:   userInfo.Avatar,
		})
	}
	return &pb.GroupInfo{
		Gid:            groupInfo.Gid,
		Name:           groupInfo.Name,
		Avatar:         groupInfo.Avatar,
		Remark:         groupInfo.Remark,
		Store:          groupInfo.Store,
		AddTime:        groupInfo.AddTime,
		Owner:          groupInfo.Owner,
		Members:        members,
		RedisDSN:       groupInfo.RedisDSN,
		RedisPassword:  groupInfo.RedisPassword,
		RedisKeyPrefix: groupInfo.RedisKeyPrefix,
		MysqlDSN:       groupInfo.MysqlDSN,
		MongodbDSN:     groupInfo.MongodbDSN,
	}
}

func (s *Service) exportTableToMysql(ctx context.Context, db *sql.DB, formatInfo *model.FormatSheet, tableName string) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	err = s.dropTable(tx, tableName)
	if err == nil {
		err = s.createTable(tx, formatInfo, tableName)
	}
	if err == nil {
		err = s.insertToTable(tx, formatInfo, tableName)
	}
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		err = tx.Rollback()
		return
	}
	return
}

func (s *Service) createTable(tx *sql.Tx, formatInfo *model.FormatSheet, tableName string) (err error) {
	createSql := "CREATE TABLE `" + tableName + "` ("
	for index, row := range formatInfo.Fields {
		fieldTy := "bigint(20)"
		if formatInfo.Types[index] == "string" {
			fieldTy = "text"
		}
		createSql += "`" + row + "` " + fieldTy + " NOT NULL COMMENT '" + formatInfo.Descs[index] + "',"
	}
	createSql += "PRIMARY KEY (`sid`) ) DEFAULT CHARSET=utf8mb4"
	_, err = tx.Exec(createSql)
	return
}

func (s *Service) insertToTable(tx *sql.Tx, formatInfo *model.FormatSheet, tableName string) (err error) {
	insertSql := "INSERT INTO `" + tableName + "` ("
	for index, field := range formatInfo.Fields {
		insertSql += "`" + field + "`"
		if index < len(formatInfo.Fields)-1 {
			insertSql += ","
		} else {
			insertSql += ")"
		}
	}
	insertSql += "VALUES("
	for rowIndex, row := range formatInfo.Content {
		for index, field := range formatInfo.Fields {
			var val = ""
			if index < len(row) {
				if tmp, ok := row[field]; ok {
					if tmpVal, ok := tmp.(string); !ok {
						val = strconv.FormatFloat(tmp.(float64), 'g', -1, 64)
					} else {
						val = tmpVal
					}
				}
			}
			insertSql += "'" + val + "'"
			if index < len(formatInfo.Fields)-1 {
				insertSql += ","
			} else {
				insertSql += ")"
			}
		}
		if rowIndex < len(formatInfo.Content)-1 {
			insertSql += ",("
		}
	}
	_, err = tx.Exec(insertSql)
	return
}

func (s *Service) dropTable(tx *sql.Tx, tableName string) (err error) {
	dropSql := "drop table if exists " + tableName
	_, err = tx.Exec(dropSql)
	if err != nil {
		return err
	}
	return
}
