package ioc

import (
	dao2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/gormx/connpool"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	promsdk "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func InitSRC(l logger.LoggerV1) SrcDB {
	return InitDB(l, "src")
}

func InitDST(l logger.LoggerV1) DstDB {
	return InitDB(l, "dst")
}

func InitDoubleWriterPool(src SrcDB, dst DstDB) *connpool.DoubleWritePool {
	pattern := viper.GetString("migrator.pattern")
	return connpool.NewDoubleWritePool(src.ConnPool, dst.ConnPool, pattern)
}

// 这个是业务用的，支持双写的 DB
func InitBizDB(pool *connpool.DoubleWritePool) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: pool,
	}))
	if err != nil {
		panic(err)
	}
	return db
}

type SrcDB *gorm.DB
type DstDB *gorm.DB

func InitDB(l logger.LoggerV1, key string) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook",
	}
	err := viper.UnmarshalKey("db." + key, &cfg)

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cb := newCallbacks(key)
	err = db.Use(cb)
	if err != nil {
		panic(err)
	}

	err = dao2.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

type Callbacks struct {
	vector *promsdk.SummaryVec
}

func (pcb *Callbacks) Name() string {
	return "prometheus-query"
}

func (pcb *Callbacks) Initialize(db *gorm.DB) error {
	pcb.registerAll(db)
	return nil
}

func newCallbacks(key string) *Callbacks {
	vector := promsdk.NewSummaryVec(promsdk.SummaryOpts{
		// 在这边，你要考虑设置各种 Namespace
		Namespace: "geekbang_daming",
		Subsystem: "webook" + key,
		Name:      "gorm_query_time",
		Help:      "统计 GORM 的执行时间",
		ConstLabels: map[string]string{
			"db": "webook",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	},
		// 如果是 JOIN 查询，table 就是 JOIN 在一起的
		// 或者 table 就是主表，A JOIN B，记录的是 A
		[]string{"type", "table"})

	pcb := &Callbacks{
		vector: vector,
	}
	promsdk.MustRegister(vector)
	return pcb
}

func (pcb *Callbacks) registerAll(db *gorm.DB) {
	// 作用于 INSERT 语句
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", pcb.after("create"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", pcb.after("update"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", pcb.after("delete"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", pcb.after("raw"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_row_after", pcb.after("row"))
	if err != nil {
		panic(err)
	}
}

func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}

func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			// 你啥都干不了
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(typ, table).
			Observe(float64(time.Since(startTime).Milliseconds()))
	}
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
