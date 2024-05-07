package ioc

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/ginx"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/gormx/connpool"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/events"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/events/fixer"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/scheduler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const topic = "migrator_interactives"

func InitFixDataConsumer(l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	client sarama.Client) *fixer.Consumer[dao.Interactive] {
	res, err := fixer.NewConsumer[dao.Interactive](client, l,
		topic, src, dst)
	if err != nil {
		panic(err)
	}
	return res
}

func InitMigratorProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer(p, topic)
}

func InitMigratorWeb(
	l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	pool *connpool.DoubleWritePool,
	producer events.Producer) *ginx.Server {
	// 在这里，有多少张表，你就初始化多少个 scheduler
	intrSch := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	engine := gin.Default()
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "Gnoloayoul",
		Subsystem: "webook_intr_admin",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})
	intrSch.RegisterRoutes(engine.Group("/migrator"))
	addr := viper.GetString("migrator.web.addr")
	return &ginx.Server{
		Addr: addr,
		Engine: engine,
	}
}