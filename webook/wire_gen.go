// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
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
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	loggerV1 := ioc.InitLogger()
	handler := jwt.NewRedisJwtHandler(cmdable)
	v := ioc.InitMiddlewares(cmdable, loggerV1, handler)
	db := ioc.InitDB(loggerV1)
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository, loggerV1)
	codeCache := cache.NewCodeRedisCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService(cmdable)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	wechatService := ioc.InitWechatService(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, handler)
	articleDAO := article.NewGORMArticleDAO(db)
	articleRepository := article2.NewArticleRepository(articleDAO, loggerV1)
	client := ioc.InitKafka()
	syncProducer := ioc.NewSyncProducer(client)
	producer := article3.NewKafkaProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	articleHandler := web.NewArticleHandler(articleService, loggerV1)
	engine := ioc.InitwebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveReadEventBatchConsumer := article3.NewInteractiveReadEventBatchConsumer(client, interactiveRepository, loggerV1)
	v2 := ioc.NewConsumers(interactiveReadEventBatchConsumer)
	app := &App{
		web:       engine,
		consumers: v2,
	}
	return app
}
