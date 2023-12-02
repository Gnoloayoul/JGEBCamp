package ioc

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web/middleware"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ginx"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ginx/middlewares/metric"
	logger2 "github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
	"time"
)

func InitwebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	(&web.ObservabilityHandler{}).RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable, l logger2.LoggerV1, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	//bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
	//	l.Debug("HTTP请求", logger2.Field{Key: "al", Value: al})
	//}).AllowReqBody(true).AllowRespBody()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	ok := viper.GetBool("web.logreq")
	//	bd.AllowReqBody(ok)
	//})
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "github-Gnoloayoul",
		Subsystem: "webook",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})
	return []gin.HandlerFunc{
		corsHdl(),
		(&metric.MiddlewareBuilder{
			Namespace:  "github-Gnoloayoul",
			Subsystem:  "webook",
			Name:       "gin_http",
			Help:       "统计 GIN 的 HTTP 接口",
			InstanceID: "my-instance-1",
		}).Build(),
		otelgin.Middleware("webook"),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login").
			Build(),
		//ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// ExposeHeaders: 用于获取需要前端截获的头，如 jwt
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
