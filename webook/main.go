package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web/middleware"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	_ "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "net"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()

	u := initUser(db)
	u.RegisterRoutes(server)

	//// 临时用的signup页面
	//server.LoadHTMLFiles("../webook-fe/signup.html")

	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(func(c *gin.Context) {
		println("这是第一个 middleware")
	})

	server.Use(func(c *gin.Context) {
   		println("这是第二个 middleware")
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr: "119.45.240.2:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1000).Build())

	server.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// ExposeHeaders: 用于获取需要前端截获的头，如 jwt
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

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

	// 使用 JWT 验证
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("acs:root278803@tcp(119.45.240.2:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}