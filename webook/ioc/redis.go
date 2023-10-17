package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "43.132.234.191:6380",
	})
	return redisClient
}

