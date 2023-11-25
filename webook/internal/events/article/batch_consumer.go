package article

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/Gnoloayoul/JGEBCamp/webook/internal/repository"
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/IBM/sarama"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventBatchConsumer(
	client sarama.Client,
	repo repository.InteractiveRepository,
	l logger.LoggerV1) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{
		client: client,
		repo: repo,
		l: l,
	}
}

func (r *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewBatchHandler[ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (r *InteractiveReadEventBatchConsumer) Consume(
	msg []*sarama.ConsumerMessage,
	ts []ReadEvent) error {
	ids := make([]int64, 0, len(ts))
	bizs := make([]string, 0, len(ts))
	for _, evt := range ts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")

	}
}

