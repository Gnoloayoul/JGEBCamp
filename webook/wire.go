//go:build wireinject

package main

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/cache"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/web"
	"github.com/Gnoloayoul/JGEBCamp/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最底层的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// dao
		dao.NewUserDAO,
		// cache
		cache.NewUserCache, cache.NewCodeRedisCache,
		// repository
		repository.NewUserRepository, repository.NewCodeRepository,
		// service
		service.NewUserService, service.NewCodeService,
		// 基于内存实现
		ioc.InitSMSService,

		web.NewUserHandler,

		ioc.InitGin, ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
