//go:build wireinject

package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/events"
	repository2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	cache2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	dao2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	service2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
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

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

var rankingServiceSet = wire.NewSet(
	repository.NewCachedRankingRepository,
	cache.NewRankingRedisCache,
	service.NewBatchRankingService,
)

func InitWebServer() *App {
	wire.Build(
		// 最底层的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		rankingServiceSet,
		interactiveSvcProvider,
		ioc.InitJobs,
		ioc.InitRankingJob,

		// consumer
		events.NewInteractiveReadEventBatchConsumer,
		article.NewKafkaProducer,

		// dao
		dao.NewUserDAO,
		article3.NewGORMArticleDAO,

		// cache
		cache.NewUserCache,
		cache.NewCodeRedisCache,

		// repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
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
