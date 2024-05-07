//go:build wireinject

package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/events"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/grpc"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/ioc"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitDST,
	ioc.InitSRC,
	ioc.InitBizDB,
	ioc.InitDoubleWriterPool,
	ioc.InitLogger,
	ioc.InitKafka,
	ioc.InitSyncProducer,
	ioc.InitRedis)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache)

var migratorProvider = wire.NewSet(
	ioc.InitMigratorWeb,
	ioc.InitFixDataConsumer,
	ioc.InitMigratorProducer)

func InitAPP() *App {
	wire.Build(interactiveSvcProvider,
		thirdProvider,
		migratorProvider,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
