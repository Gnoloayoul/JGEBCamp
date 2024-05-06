package startup

import (
	events2 "github.com/Gnoloayoul/JGEBCamp/webook/interactive/events"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/events"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{localhost + ":9094"}, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}