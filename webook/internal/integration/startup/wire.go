//go:build wireinject

package startup

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/Gnoloayoul/JGEBCamp/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,

		cache.NewCodeRedisCache,
		dao.NewGORMArticleDAO,
		repository.NewCodeRepository,
		repository.NewArticleRepository,

		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSMSService,

		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		service.NewCodeService,
		service.NewArticleService,

		// handler 部分
		web.NewUserHandler,
		web.NewWechatHandler,
		web.NewArticleHandler,
		InitWechatHandlerConfig,
		ijwt.NewRedisJwtHandler,

		// gin 的中间件
		ioc.InitMiddlewares,

		// Web 服务器
		ioc.InitwebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		dao.NewGORMArticleDAO,
		service.NewArticleService,
		web.NewArticleHandler,
		repository.NewArticleRepository,
	)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJwtHandler)
	return ijwt.NewRedisJwtHandler(nil)
}
