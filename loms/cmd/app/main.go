package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"route256/libs/basemetrics"
	"route256/libs/healthcheck"
	"route256/libs/logger"
	"route256/libs/serverwrapper"
	"route256/libs/tracing"
	"route256/libs/transactor"
	lomsV1 "route256/loms/internal/api/loms_v1"
	"route256/loms/internal/client/db"
	"route256/loms/internal/client/kafka"
	"route256/loms/internal/config"
	"route256/loms/internal/cron"
	"route256/loms/internal/cron/jobs/cancelreservationduetimeout"
	"route256/loms/internal/cron/jobs/readoutboxsend"
	"route256/loms/internal/domain"
	orderStatusChanges "route256/loms/internal/producer/order_status_changes"
	"route256/loms/internal/repository"
	desc "route256/loms/pkg/loms_v1"
)

var (
	develMode = flag.Bool("devel", false, "development mode")
)

func main() {
	flag.Parse()
	err := config.New()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	logger.Init(*develMode, config.Data.LoggerLevel)
	tracing.Init(config.Data.ServiceName)

	ctx := context.Background()

	dbClient, err := db.New(ctx, &config.Data.DbConfig)
	if err != nil {
		logger.Fatalf("failed to init database: %v", err)
	}

	tm := transactor.NewTransactor(dbClient)
	lr := repository.NewLomsRepository(tm)

	sp, err := kafka.NewSyncProducer(config.Data.Kafka.Brokers)
	if err != nil {
		logger.Fatalf("failed to create sync producer: %v", err)
	}
	pr := orderStatusChanges.NewProducer(sp, config.Data.CronJobs.ReadOutBoxSendJob.TopicToSend)

	rs := readoutboxsend.NewJob(lr, pr)
	cr := cancelreservationduetimeout.NewJob(lr, tm)
	jm := cron.NewJobsManager(cr, rs)

	err = jm.StartAllJobs(ctx)
	if err != nil {
		logger.Fatalf("can't start job: %v", err)
	}

	bl := domain.NewBuisnessLogic(lr, tm)

	grpcServer := serverwrapper.NewGrpcServer(config.Data.GrpcPort)

	desc.RegisterLomsV1Server(grpcServer.GetServer(), lomsV1.NewLomsV1(bl))

	go func() {
		logger.Infof("serving metrics at :%v", config.Data.MetricsPort)

		http.HandleFunc("/healthcheck", healthcheck.Handler)
		http.Handle("/metrics", basemetrics.New())
		err = http.ListenAndServe(fmt.Sprintf(":%d", config.Data.MetricsPort), nil)
		if err != nil {
			logger.Fatalf("failed to start serving metrics: %v", err)
		}
	}()

	logger.Infof("server listening at %v", grpcServer.GetListener().Addr())
	if err = grpcServer.Serve(); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
