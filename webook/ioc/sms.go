package ioc

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service{
	// 这里随便换
	return memory.NewService()
}
