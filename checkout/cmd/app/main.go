package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	checkoutV1 "route256/checkout/internal/api/checkout_v1"
	"route256/checkout/internal/clients/db"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/ps"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/basemetrics"
	"route256/libs/cache"
	"route256/libs/healthcheck"
	"route256/libs/logger"
	"route256/libs/serverwrapper"
	"route256/libs/tracing"
	"route256/libs/transactor"
)

var (
	develMode = flag.Bool("devel", false, "development mode")
)

func main() {
	err := config.New()
	if err != nil {
		log.Fatal("config init", err)
	}

	logger.Init(*develMode, config.Data.LoggerLevel)
	tracing.Init(config.Data.ServiceName)

	ctx := context.Background()

	lomsClient := loms.NewClient(ctx)
	defer lomsClient.Conn.Close()

	psClient := ps.NewClient(ctx)
	defer psClient.Conn.Close()

	cacheSettings := config.Data.Cache
	ch := cache.NewCache[string](ctx, cacheSettings.Ttl, cacheSettings.Buckets, cacheSettings.LruCapacity)

	dbClient, err := db.New(context.Background(), &config.Data.DbConfig)
	if err != nil {
		logger.Fatalf("failed to init database: %v", err)
	}

	tm := transactor.NewTransactor(dbClient)

	cr := repository.NewCheckoutRepository(tm)

	businessLogic := domain.NewBuisnessLogic(lomsClient, psClient, cr, tm, ch)

	grpcServer := serverwrapper.NewGrpcServer(config.Data.GrpcPort)
	desc.RegisterCheckoutV1Server(grpcServer.GetServer(), checkoutV1.NewCheckoutV1(businessLogic))

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
