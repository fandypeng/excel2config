package helper

import (
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/container/pool"
	"github.com/go-kratos/kratos/pkg/database/sql"
	xtime "github.com/go-kratos/kratos/pkg/time"
	"time"
)

func ParseRedisDsn(dsn, pwd string) (conf *redis.Config) {
	return &redis.Config{
		Config: &pool.Config{
			Active:      100,
			Idle:        10,
			IdleTimeout: 10,
			WaitTimeout: 0,
			Wait:        false,
		},
		Name:         "demo",
		Proto:        "tcp",
		Addr:         dsn,
		Auth:         pwd,
		DialTimeout:  xtime.Duration(time.Second * 1),
		ReadTimeout:  xtime.Duration(time.Second * 1),
		WriteTimeout: xtime.Duration(time.Second * 1),
		SlowLog:      xtime.Duration(time.Second * 5),
	}
}

func ParseMysqlDsn(dsn string) (conf *sql.Config) {
	return &sql.Config{
		DSN:          dsn,
		ReadDSN:      []string{},
		Active:       100,
		Idle:         10,
		IdleTimeout:  xtime.Duration(time.Second * 1),
		QueryTimeout: xtime.Duration(time.Second * 1),
		ExecTimeout:  xtime.Duration(time.Second * 1),
		TranTimeout:  xtime.Duration(time.Second * 1),
		Breaker:      nil,
	}
}
