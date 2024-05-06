package startup

import "github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
