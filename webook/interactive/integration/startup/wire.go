//go:build wireinject

package startup

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/grpc"
	repository2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	grpcRepo "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/grpc"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	cache2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	dao2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	service2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var interactiveRepoProvider = wire.NewSet(
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
	repository2.NewCachedInteractiveRepository)
var interactiveSvcProvider = wire.NewSet(
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

func InitInteractiveRepoService() repository.InteractiveRepository {
	wire.Build(thirdProvider, interactiveRepoProvider)
	return repository.NewCachedInteractiveRepository(nil, nil, nil)
}

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}

func InitInteractiveRepoGRPCServer() *grpcRepo.InteractiveGRPCRepositoryServer {
	wire.Build(thirdProvider, interactiveRepoProvider, grpcRepo.NewInteractiveGRPCRepositoryServer)
	return new(grpcRepo.InteractiveGRPCRepositoryServer)
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
