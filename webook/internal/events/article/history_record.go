package article

import (
	"context"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/saramax"
	"github.com/IBM/sarama"
)

type HistoryReadEventConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
}

func NewHistoryReadEventConsumer(
	client sarama.Client,
	l logger.LoggerV1) *HistoryReadEventConsumer {
	return &HistoryReadEventConsumer{
		client: client,
		l:      l,
	}
}

func (r *HistoryReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (r *HistoryReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//return r.repo.AddRecord(ctx, t.Aid, t.Uid)
	panic("implement me")
}
