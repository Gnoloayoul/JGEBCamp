package startup

import (
	"github.com/Gnoloayoul/JGEBCamp/webook_homework/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
