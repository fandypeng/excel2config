package dao

import (
	"context"
	pb "excel2config/api"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/go-kratos/kratos/pkg/time"

	"github.com/google/wire"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, NewMC, NewMongo)

//go:generate kratos tool genbts
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	// bts: -nullcache=&model.Article{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Article(c context.Context, id int64) (*model.Article, error)

	ExcelList(c context.Context, lastTime, limit int64, groupId string) (list []*pb.SimpleExcel, err error)
	CreateExcel(c context.Context, uid, gridKey, remark, groupId string) (eid string, err error)
	UpdateExcel(c context.Context, id, remark string, contributers []string) (err error)
	DeleteExcel(c context.Context, id, gridKey string) (err error)
	ExcelInfo(c context.Context, gridKey string) (excelInfo *model.Excel, err error)

	LoadExcel(c context.Context, gridKey string) (sheets []*model.Sheet, err error)
	LoadExcelSheet(c context.Context, gridKey string, indexs []string) (sheets map[string][]model.Cell, err error)
	LoadAllSheet(c context.Context, gridKey string) (sheets []*model.Sheet, err error)
	LoadSheetByName(c context.Context, gridKey, sheetName string) (sheet *model.Sheet, err error)
	UpdateGridValue(c context.Context, gridKey string, req *model.UpdateV) (err error)
	UpdateGridMulti(c context.Context, gridKey string, req *model.UpdateRV) (err error)
	UpdateGridConfig(c context.Context, gridKey string, req *model.UpdateCG) (err error)
	UpdateGridCommon(c context.Context, gridKey string, req *model.UpdateCommon) (err error)
	UpdateCalcChain(c context.Context, gridKey string, req *model.UpdateCalcChain) (err error)
	UpdateRowColumn(c context.Context, gridKey string, req *model.UpdateRowColumn) (err error)
	UpdateFilter(c context.Context, gridKey string, req *model.UpdateFilter) (err error)

	AddSheet(c context.Context, gridKey string, req *model.AddSheet) (err error)
	CopySheet(c context.Context, gridKey string, req *model.CopySheet) (err error)
	DeleteSheet(c context.Context, gridKey string, req *model.DeleteSheet) (err error)
	RecoverSheet(c context.Context, gridKey string, req *model.RecoverSheet) (err error)
	UpdateSheetOrder(c context.Context, gridKey string, req *model.UpdateSheetOrder) (err error)
	ToggleSheet(c context.Context, gridKey string, req *model.ToggleSheet) (err error)
	HideOrShowSheet(c context.Context, gridKey string, req *model.HideOrShowSheet) (err error)

	GetUser(c context.Context, email string) (userInfo *model.UserInfo, err error)
	GetUserInfos(c context.Context, uids []string) (userInfos map[string]*model.UserInfo, err error)
	GetUserByUid(c context.Context, uid string) (userInfo *model.UserInfo, err error)
	AddUser(c context.Context, email, passwd, name string, loginType int) (userInfo *model.UserInfo, err error)
	GetToken(c context.Context, uid string) (token string, err error)
	SaveToken(c context.Context, uid, token string) (err error)
	GetUserInfosByKeyword(c context.Context, keyword string) (userInfos map[string]*model.SimpleUserInfo, err error)

	GroupList(c context.Context, uid string) (groupList []*model.GroupInfo, err error)
	GroupAdd(c context.Context, groupInfo *model.GroupInfo) (groupId string, err error)
	GroupUpdate(c context.Context, groupInfo *model.GroupInfo) (err error)
	GroupInfo(c context.Context, groupId string) (groupInfo *model.GroupInfo, err error)
}

// dao dao.
type dao struct {
	db         *sql.DB
	redis      *redis.Redis
	mc         *memcache.Memcache
	mongo      *mongo.Client
	cache      *fanout.Fanout
	demoExpire int32
}

// New new a dao and return.
func New(r *redis.Redis, mc *memcache.Memcache, db *sql.DB, mongo *mongo.Client) (d Dao, cf func(), err error) {
	return newDao(r, mc, db, mongo)
}

func newDao(r *redis.Redis, mc *memcache.Memcache, db *sql.DB, mongo *mongo.Client) (d *dao, cf func(), err error) {
	var cfg struct {
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db:         db,
		redis:      r,
		mc:         mc,
		mongo:      mongo,
		cache:      fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	d.db.Ping(ctx)
	return nil
}
