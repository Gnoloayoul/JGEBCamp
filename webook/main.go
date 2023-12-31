package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Gnoloayoul/JGEBCamp/webook/ioc"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//
	//rdb := initRedis()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)
	initViperV1()
	initLogger()

	closeFunc := ioc.InitOTEL()
	initPrometheus()
	keys := viper.AllKeys()
	println(keys)

	app := InitWebServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	//// 临时用的signup页面
	//server.LoadHTMLFiles("../webook-fe/signup.html")

	//server := gin.Default()
	app.cron.Start()

	server := app.web
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello here is k8s")
	})
	server.Run(":8080")
	// 一分钟内你要关完，要退出
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	closeFunc(ctx)

	ctx = app.cron.Stop()
	// 想办法 close ？？
	// 这边可以考虑超时强制退出，防止有些任务，执行特别长的时间
	tm := time.NewTimer(time.Minute * 10)
	select {
	case <-tm.C:
	case <-ctx.Done():
	}
	// 作业
	//server.Run(":8081")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(func(c *gin.Context) {
		println("这是第一个 middleware")
	})

	server.Use(func(c *gin.Context) {
		println("这是第二个 middleware")
	})

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: "146.56.252.134:6380",
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1000).Build())

	// step1

	// old
	//store := cookie.NewStore([]byte("secret"))

	// new 单机
	//store := memstore.NewStore([]byte("h7oUXRzcGPyJbZJfq68iGChnzA0iJBfJ"),
	//	[]byte("aRNaEVNTV5IOzXbatCQuQCkwNteyJwPe"))

	// new 多机 Redis
	//store, err := redis.NewStore(16, "tcp", "119.45.240.2:6379", "",
	//	// authentication key, encryption key
	//	[]byte("h7oUXRzcGPyJbZJfq68iGChnzA0iJBfJ"),
	//	[]byte("aRNaEVNTV5IOzXbatCQuQCkwNteyJwPe"))
	//if err != nil {
	//	panic(err)
	//}
	store := memstore.NewStore(
		[]byte("h7oUXRzcGPyJbZJfq68iGChnzA0iJBfJ"),
		[]byte("aRNaEVNTV5IOzXbatCQuQCkwNteyJwPe"))

	server.Use(sessions.Sessions("mysession", store))

	//// step3
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())

	return server
}

func initViper() {
	viper.SetDefault("db.mysql.dsn",
		"root:root@tcp(localhost:3306)/mysql")
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	// 现实中，有很多格式，JSON，XML，YAML，TOML，ini
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")
	//viper.AddConfigPath("/tmp/config")
	//viper.AddConfigPath("/etc/webook")
	// 读取配置到 viper 里面，或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}

func initViperV1() {
	cfile := pflag.String("config",
		"config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:3306)/mysql")
	//viper.SetConfigFile("config/dev.yaml")
	//viper.KeyDelimiter("-")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果你不 replace，直接用 zap.L()，你啥都打不出来。
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")

	type Demo struct {
		Name string `json:"name"`
	}
	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:6379"
`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		// 通过 webook 和其他使用 etcd 的区别出来
		"http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}
