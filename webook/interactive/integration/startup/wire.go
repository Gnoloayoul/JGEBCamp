//go:build wireinject

package startup

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/grpc"
	repository2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository"
	cache2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/cache"
	dao2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	service2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
