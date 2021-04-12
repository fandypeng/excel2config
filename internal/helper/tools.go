package helper

import (
	"encoding/json"
	"excel2config/internal/model"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/container/pool"
	"github.com/go-kratos/kratos/pkg/database/sql"
	xtime "github.com/go-kratos/kratos/pkg/time"
	"math/rand"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789=")

func GenerateRandomStr(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	return string(b)
}

func GetIntroduceCellData() []model.Cell {
	jsonstr := `[{"r":0,"v":{"mc":{"cs":6,"rs":6},"ct":{"s":[{"v":"配置表头分为三行\r\n1、第一行是字段名，用英文表示\r\n2、第二行是字段类型，int代表数值类型，string代表字符串类型\r\n3、第三行是字段备注，中英文均可，建议不超过30个字\r\n\r\n首次建表请点击下方的+号，新建sheet，按照上面的格式配置数据。"}],"t":"inlineStr","fa":"General"}},"c":0},{"c":0,"r":1,"v":{"mc":{}}},{"c":0,"r":2,"v":{"mc":{}}},{"c":0,"r":3,"v":{"mc":{}}},{"c":0,"r":4,"v":{"mc":{}}},{"c":0,"r":5,"v":{"mc":{}}},{"c":1,"r":0,"v":{"mc":{}}},{"c":1,"r":1,"v":{"mc":{}}},{"c":1,"r":2,"v":{"mc":{}}},{"c":1,"r":3,"v":{"mc":{}}},{"c":1,"r":4,"v":{"mc":{}}},{"c":1,"r":5,"v":{"mc":{}}},{"v":{"mc":{}},"c":2,"r":0},{"r":1,"v":{"mc":{}},"c":2},{"v":{"mc":{}},"c":2,"r":2},{"v":{"mc":{}},"c":2,"r":3},{"v":{"mc":{}},"c":2,"r":4},{"c":2,"r":5,"v":{"mc":{}}},{"c":3,"r":0,"v":{"mc":{}}},{"r":1,"v":{"mc":{}},"c":3},{"c":3,"r":2,"v":{"mc":{}}},{"c":3,"r":3,"v":{"mc":{}}},{"c":3,"r":4,"v":{"mc":{}}},{"c":3,"r":5,"v":{"mc":{}}},{"c":4,"r":0,"v":{"mc":{}}},{"c":4,"r":1,"v":{"mc":{}}},{"c":4,"r":2,"v":{"mc":{}}},{"c":4,"r":3,"v":{"mc":{}}},{"c":4,"r":4,"v":{"mc":{}}},{"r":5,"v":{"mc":{}},"c":4},{"c":5,"r":0,"v":{"mc":{}}},{"r":1,"v":{"mc":{}},"c":5},{"c":5,"r":2,"v":{"mc":{}}},{"c":5,"r":3,"v":{"mc":{}}},{"c":5,"r":4,"v":{"mc":{}}},{"c":5,"r":5,"v":{"mc":{}}}]`
	cells := make([]model.Cell, 0)
	json.Unmarshal([]byte(jsonstr), &cells)
	return cells
}

func GetIntroduceSheetConfig() *model.SheetConfig {
	jsonstr := `{"merge":{"0_0":{"rs":6,"cs":6,"r":0,"c":0}},"rowlen":{"0":95.0898447036743}}`
	conf := new(model.SheetConfig)
	json.Unmarshal([]byte(jsonstr), conf)
	return conf
}

func GetSheetInitCellData() []model.Cell {
	jsonstr := `[{"v":{"ct":{"fa":"General","t":"g"},"m":"sid","v":"sid"},"c":0,"r":0},{"c":0,"r":1,"v":{"v":"int","ct":{"t":"g","fa":"General"},"m":"int"}},{"v":{"m":"流水ID","v":"流水ID","ct":{"fa":"General","t":"g"}},"c":0,"r":2}]`
	cells := make([]model.Cell, 0)
	json.Unmarshal([]byte(jsonstr), &cells)
	return cells
}
