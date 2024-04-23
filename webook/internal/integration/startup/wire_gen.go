// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	repository2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	cache2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	dao2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	service2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
	article3 "github.com/Gnoloayoul/JGEBCamp/webook/internal/events/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	article2 "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/Gnoloayoul/JGEBCamp/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	loggerV1 := InitLog()
	handler := jwt.NewRedisJwtHandler(cmdable)
	v := ioc.InitMiddlewares(cmdable, loggerV1, handler)
	gormDB := InitTestDB()
	userDAO := dao.NewUserDAO(gormDB)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository, loggerV1)
	codeCache := cache.NewCodeRedisCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService(cmdable)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := InitPhantomWechatService(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	articleDAO := article.NewGORMArticleDAO(gormDB)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := article2.NewArticleRepository(articleDAO, articleCache, userRepository, loggerV1)
	client := InitKafka()
	syncProducer := NewSyncProducer(client)
	producer := article3.NewKafkaProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	interactiveDAO := dao2.NewGORMInteractiveDAO(gormDB)
	interactiveCache := cache2.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository2.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service2.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceClient := ioc.InitGRPCInteractiveServiceClient(interactiveService)
	articleHandler := web.NewArticleHandler(articleService, interactiveServiceClient, loggerV1)
	engine := ioc.InitwebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	return engine
}

func InitArticleHandler(dao3 article.ArticleDAO) *web.ArticleHandler {
	cmdable := InitRedis()
	articleCache := cache.NewRedisArticleCache(cmdable)
	gormDB := InitTestDB()
	userDAO := dao.NewUserDAO(gormDB)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	loggerV1 := InitLog()
	articleRepository := article2.NewArticleRepository(dao3, articleCache, userRepository, loggerV1)
	client := InitKafka()
	syncProducer := NewSyncProducer(client)
	producer := article3.NewKafkaProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	interactiveDAO := dao2.NewGORMInteractiveDAO(gormDB)
	interactiveCache := cache2.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository2.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service2.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceClient := ioc.InitGRPCInteractiveServiceClient(interactiveService)
	articleHandler := web.NewArticleHandler(articleService, interactiveServiceClient, loggerV1)
	return articleHandler
}

func InitUserSvc() service.UserService {
	gormDB := InitTestDB()
	userDAO := dao.NewUserDAO(gormDB)
	cmdable := InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	loggerV1 := InitLog()
	userService := service.NewUserService(userRepository, loggerV1)
	return userService
}

func InitJwtHdl() jwt.Handler {
	cmdable := InitRedis()
	handler := jwt.NewRedisJwtHandler(cmdable)
	return handler
}

// wire.go:

var thirdProvider = wire.NewSet(InitRedis,
	NewSyncProducer,
	InitKafka,
	InitTestDB, InitLog)

var userSvcProvider = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewUserRepository, service.NewUserService)

var interactiveSvcProvider = wire.NewSet(service2.NewInteractiveService, repository2.NewCachedInteractiveRepository, dao2.NewGORMInteractiveDAO, cache2.NewRedisInteractiveCache)

var articleSvcProvider = wire.NewSet(article.NewGORMArticleDAO, cache.NewRedisArticleCache, article2.NewArticleRepository, service.NewArticleService)
