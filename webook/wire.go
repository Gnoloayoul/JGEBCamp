//go:build wireinject

package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/events/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	article2 "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	article3 "github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao/article"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	ijwt "github.com/Gnoloayoul/JGEBCamp/webook/internal/web/jwt"
	"github.com/Gnoloayoul/JGEBCamp/webook/ioc"
	"github.com/google/wire"
)

func InitWebServer() *App {
	wire.Build(
		// 最底层的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		// consumer
		article.NewInteractiveReadEventBatchConsumer,
		article.NewKafkaProducer,

		// dao
		dao.NewUserDAO,
		article3.NewGORMArticleDAO,
		dao.NewGORMInteractiveDAO,

		// cache
		cache.NewRedisInteractiveCache,
		cache.NewUserCache,
		cache.NewCodeRedisCache,


		// repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewCachedInteractiveRepository,
		article2.NewArticleRepository,

		// service
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// 基于内存实现
		ioc.InitSMSService,
		ioc.InitWechatService,

		//ioc.NewWechatHandlerconfig,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,

		ijwt.NewRedisJwtHandler,

		ioc.InitwebServer,
		ioc.InitMiddlewares,

		// 组装我这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
