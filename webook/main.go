package main

import (
	_ "github.com/Gnoloayoul/JGEBCamp/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/sessions"
	_ "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	_ "net"
	"net/http"
)

func main() {
	//db := initDB()
	//server := initWebServer()
	//
	//rdb := initRedis()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)
	initViper()

	server := InitWebServer()

	//// 临时用的signup页面
	//server.LoadHTMLFiles("../webook-fe/signup.html")

	//server := gin.Default()
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello here is k8s")
	})
	server.Run(":8081")
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

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("before Hello")
	zap.ReplaceGlobals(logger)
	zap.L().Info("after Hello")
}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
