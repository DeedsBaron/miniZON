package producer

import "context"

type KafkaProducer interface {
	SendMessage(ctx context.Context, key string, message []byte) error
}
