package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"route256/libs/basemetrics"
)

func MetricsInterceptor(ctx context.Context,
	method string, req interface{},
	reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	timeStart := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	elapsed := time.Since(timeStart)
	basemetrics.ClientHistogramResponseTime.
		WithLabelValues(cc.Target()).
		Observe(elapsed.Seconds())
	return err
}
