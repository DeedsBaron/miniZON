package order_status_changes

import (
	"context"
	"time"

	"route256/libs/logger"

	"github.com/Shopify/sarama"
)

type OrderStatusChangesProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(producer sarama.SyncProducer, topic string) *OrderStatusChangesProducer {
	return &OrderStatusChangesProducer{
		producer: producer,
		topic:    topic,
	}
}

func (o *OrderStatusChangesProducer) SendMessage(ctx context.Context, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     o.topic,
		Partition: -1,
		Value:     sarama.StringEncoder(value),
		Key:       sarama.StringEncoder(key),
		Timestamp: time.Now(),
	}

	partition, offset, err := o.producer.SendMessage(msg)
	if err != nil {
		logger.Errorf("sending in kafka, error occurred, %v", err)
		return err
	}

	logger.Infow("sent to kafa",
		"topic", o.topic,
		"partition", partition,
		"offset", offset,
		"key", key,
		"value", string(value))
	return nil
}
