package saramax

import (
	"github.com/Gnoloayoul/JGEBCamp/webook/pkg/logger"
	"github.com/IBM/sarama"
)

type HandlerV1[T any] func(msg *sarama.ConsumerMessage, t T) error

type Handler[T any] struct {

}

func NewHandler[T any](l logger.LoggerV1, fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {

}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {}

func (h Handler[T]) Setup(session sarama.ConsumerGroupSession) error {}