package ioc

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/events"
	"github.com/Gnoloayoul/JGEBCamp/webook/interactive/repository/dao"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/migrator/events/fixer"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	//// 云服务器访问操作
	//cloudServerIP := cfg.Addrs[:len(cfg.Addrs) - 5] // 替换为实际的云服务器IP地址
	//
	//// 将云服务器的IP地址写入Kafka的配置
	//viper.Set("kafka.advertised.listeners", fmt.Sprintf("PLAINTEXT://%s:9092,EXTERNAL://%s:9094", cloudServerIP, cloudServerIP))
	//
	//// 保存配置文件
	//err = viper.WriteConfig()
	//if err != nil {
	//	panic(fmt.Errorf("Failed to write config file: %w", err))
	//}

	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

// 规避 wire 的问题
type fixerInteractive *fixer.Consumer[dao.Interactive]


// NewConsumers 面临的问题依旧是所有的 Consumer 在这里注册一下
func NewConsumers(intr *events.InteractiveReadEventConsumer,
	fix *fixer.Consumer[dao.Interactive]) []saramax.Consumer {
	return []saramax.Consumer{
		intr,
		fix,
	}
}
