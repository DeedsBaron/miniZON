package serverinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc/status"
	"route256/libs/basemetrics"

	"google.golang.org/grpc"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	basemetrics.ServerRequestsCounter.WithLabelValues(info.FullMethod).Inc()

	timeStart := time.Now()

	res, err := handler(ctx, req)

	statusCode := status.Code(err)

	elapsed := time.Since(timeStart)

	basemetrics.ServerHistogramResponseTime.
		WithLabelValues(info.FullMethod, statusCode.String()).
		Observe(elapsed.Seconds())
	return res, err
}
