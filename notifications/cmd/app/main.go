package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"route256/libs/healthcheck"
	"route256/notifications/internal/client/kafka"
	"route256/notifications/internal/config"
	"route256/notifications/internal/consumer"
	"route256/notifications/internal/logger"
)

var (
	develMode = flag.Bool("devel", false, "development mode")
)

func main() {
	err := config.New()
	if err != nil {
		log.Fatalf("failed to config init: %v", err)
	}
	logger.Init(*develMode)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		logger.Infof("serving metrics at :%v", config.Data.MetricsPort)

		http.HandleFunc("/healthcheck", healthcheck.Handler)

		err = http.ListenAndServe(fmt.Sprintf(":%d", config.Data.MetricsPort), nil)
		if err != nil {
			logger.Fatalf("failed to start serving metrics: %v", err)
		}
	}()

	cg := consumer.NewConsumerGroup(kafka.NewConsumer())
	cg.StartConsumingCycle(ctx)
	cancel()
}
