package startup

import (
	"github.com/Gnoloayoul/JGEBCamp/webook_homework/internal/service/oauth2/wechat"
	"github.com/Gnoloayoul/JGEBCamp/webook_homework/pkg/logger"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
