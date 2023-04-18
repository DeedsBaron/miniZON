package serverinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"route256/libs/logger"
)

const dateLayout = "2006-01-02"

// LoggingInterceptor ...
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Infow("handler request",
		"handler", info.FullMethod,
		"req", req)

	res, err := handler(ctx, req)
	if err != nil {
		logger.Errorw("handler error",
			"handler", info.FullMethod,
			"err", err.Error())
		return nil, err
	}

	logger.Infow("handler response",
		"handler", info.FullMethod,
		"resp", res)

	return res, nil
}
