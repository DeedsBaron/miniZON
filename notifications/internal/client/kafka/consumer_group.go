package kafka

import (
	"github.com/Shopify/sarama"
	"route256/notifications/internal/logger"
)

type Consumer interface {
	Ready() <-chan bool
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

// Consumer represents a Sarama consumer group consumer
type consumer struct {
	ready chan bool
}

func NewConsumer() *consumer {
	return &consumer{
		ready: make(chan bool),
	}
}

func (consumer *consumer) Ready() <-chan bool {
	return consumer.ready
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			logger.Infow("message claimed",
				"topic", message.Topic,
				"offset", message.Offset,
				"partition", message.Partition,
				"key", string(message.Key),
				"value", string(message.Value),
				"msgTimestamp", message.Timestamp)
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
