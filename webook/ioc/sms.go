package ioc

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 这里随便换
	return memory.NewService()
}
