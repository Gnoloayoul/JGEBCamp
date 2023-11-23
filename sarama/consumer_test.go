package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()

}

type testConsumerGroupHandler struct {
}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {

}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {

}

func (t testConsumerGroupHandler) ConsumeClaim(

	) error {

}

func (t testConsumerGroupHandler) ConsumeClaimV1(

) error {

}

type MyBizMsg struct {
	Name string
}

// 返回只读的 channel
func ChannelV1() <-chan struct{} {
	panic("implement me")
}

// 返回可读可写的 channel
func ChannelV2() chan struct{} {
	panic("implement me")
}

// 返回只写的 channel
func ChannelV3() chan<- struct{} {
	panic("implement me")
}