package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "read_article",
			Value: sarama.StringEncoder(`{"aid": 1, "uid": 123}`),
		})
		assert.NoError(t, err)
	}
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewAsyncProducer(addrs, cfg)
	require.NoError(t, err)
	msgCh := producer.Input()

	go func() {
		for {
			msg := &sarama.ProducerMessage{
				Topic: "test_topic",
				Key: sarama.StringEncoder("oid-123"),
				Value: sarama.StringEncoder("Hello, 这是一条信息 A"),
				Headers: []sarama.RecordHeader{
					{
						Key: []byte("trace_id"),
						Value: []byte("123465"),
					},
				},
				Metadata: "这是metadata",
			}
			select {
			case msgCh <- msg:
				// default:
			}
		}
	}()

	errCh := producer.Errors()
	surrCh := producer.Successes()

	for {
		select {
		case err := <-errCh:
			t.Log("发送出来问题", err.Err)
		case <-surrCh:
			t.Log("发送成功")
		}
	}
}

type JSONEncoder struct {
	Data any
}