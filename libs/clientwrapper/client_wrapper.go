package clientwrapper

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256/libs/clientinterceptors"
	"route256/libs/logger"
)

func NewGrpcConnection(ctx context.Context, dsn string) *grpc.ClientConn {
	connectionCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	conn, err := grpc.DialContext(connectionCtx, dsn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(clientinterceptors.MetricsInterceptor,
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		logger.Fatalf("failed to create GRPC connection: %v", err)
	}
	return conn
}
