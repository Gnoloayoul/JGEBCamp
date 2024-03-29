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
	db := ioc.InitDB(loggerV1)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	server := ioc.InitGRPCxServer(interactiveServiceServer)
	client := ioc.InitKafka()
	interactiveReadEventBatchConsumer := events.NewInteractiveReadEventBatchConsumer(client, interactiveRepository, loggerV1)
	v := ioc.NewConsumers(interactiveReadEventBatchConsumer)
	app := &App{
		server:    server,
		consumers: v,
	}
	return app
}

// wire.go:

var thirdProvider = wire.NewSet(ioc.InitRedis, ioc.InitDB, ioc.InitLogger, ioc.InitKafka)

var interactiveSvcProvider = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)
