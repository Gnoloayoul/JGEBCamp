// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func InitAPP() *App {
	loggerV1 := ioc.InitLogger()
	srcDB := ioc.InitSRC(loggerV1)
	dstDB := ioc.InitDST(loggerV1)
	doubleWritePool := ioc.InitDoubleWriterPool(srcDB, dstDB)
	db := ioc.InitBizDB(doubleWritePool)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	server := ioc.InitGRPCxServer(interactiveServiceServer)
	client := ioc.InitKafka()
	interactiveReadEventConsumer := events.NewInteractiveReadEventConsumer(client, loggerV1, interactiveRepository)
	consumer := ioc.InitFixDataConsumer(loggerV1, srcDB, dstDB, client)
	v := ioc.NewConsumers(interactiveReadEventConsumer, consumer)
	syncProducer := ioc.InitSyncProducer(client)
	producer := ioc.InitMigratorProducer(syncProducer)
	ginxServer := ioc.InitMigratorWeb(loggerV1, srcDB, dstDB, doubleWritePool, producer)
	app := &App{
		server:    server,
		consumers: v,
		webAdmin:  ginxServer,
	}
	return app
}

// wire.go:

var thirdProvider = wire.NewSet(ioc.InitDST, ioc.InitSRC, ioc.InitBizDB, ioc.InitDoubleWriterPool, ioc.InitLogger, ioc.InitKafka, ioc.InitSyncProducer, ioc.InitRedis)

var interactiveSvcProvider = wire.NewSet(service.NewInteractiveService, repository.NewCachedInteractiveRepository, dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache)

var migratorProvider = wire.NewSet(ioc.InitMigratorWeb, ioc.InitFixDataConsumer, ioc.InitMigratorProducer)
