package consumer

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"route256/notifications/internal/client/kafka"
	"route256/notifications/internal/config"
	"route256/notifications/internal/logger"
)

type Group struct {
	client   sarama.ConsumerGroup
	consumer kafka.Consumer
}

func NewConsumerGroup(consumer kafka.Consumer) *Group {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.MaxVersion
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(config.Data.Kafka.Brokers, config.Data.Kafka.ConsumerGroup, cfg)
	if err != nil {
		logger.Fatalf("Error creating consumer group client: %v", err)
	}
	logger.Infof("Counsumer group client was successfully created")
	return &Group{
		client:   client,
		consumer: consumer,
	}
}

func (g *Group) StartConsumingCycle(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// serverwrapper-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := g.client.Consume(ctx, []string{config.Data.Kafka.Topic}, g.consumer); err != nil {
				logger.Fatalf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()
	<-g.consumer.Ready()
	logger.Infof("Consumer successfully subscribed to topic: %s\n", config.Data.Kafka.Topic)
	wg.Wait()
	if err := g.client.Close(); err != nil {
		logger.Fatalf("Error closing client: %v", err)
	}
}
