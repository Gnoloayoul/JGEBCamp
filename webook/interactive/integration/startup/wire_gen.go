// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/grpc"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitInteractiveService() service.InteractiveService {
	gormDB := InitTestDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	loggerV1 := InitLog()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	return interactiveService
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	gormDB := InitTestDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	loggerV1 := InitLog()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	return interactiveServiceServer
}

// wire.go:

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)

var interactiveSvcProvider = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)
