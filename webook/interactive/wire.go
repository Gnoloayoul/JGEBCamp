//go:build wireinject

package main

var thirdProvider = wire.NewSet(InitRedis, InitDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService)

