package dao

import (
	"context"
	"fmt"

	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/cache/memcache"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

//go:generate kratos tool genmc
type _mc interface {
	// mc: -key=keyArt -type=get
	CacheArticle(c context.Context, id int64) (*model.Article, error)
	// mc: -key=keyArt -expire=d.demoExpire
	AddCacheArticle(c context.Context, id int64, art *model.Article) (err error)
	// mc: -key=keyArt
	DeleteArticleCache(c context.Context, id int64) (err error)
}

func NewMC() (mc *memcache.Memcache, cf func(), err error) {
	var (
		cfg memcache.Config
		ct paladin.TOML
	)
	if err = paladin.Get("memcache.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	mc =  memcache.New(&cfg)
	cf = func() {mc.Close()}
	return
}

func (d *dao) PingMC(ctx context.Context) (err error) {
	if err = d.mc.Set(ctx, &memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

func keyArt(id int64) string {
	return fmt.Sprintf("art_%d", id)
}
