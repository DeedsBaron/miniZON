package serverwrapper

import (
	"fmt"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"route256/libs/logger"
	"route256/libs/serverinterceptors"
)

type grpcServer struct {
	lis        net.Listener
	grpcServer *grpc.Server
}

func NewGrpcServer(port int, interceptors ...grpc.UnaryServerInterceptor) *grpcServer {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	resultInterceptors := []grpc.UnaryServerInterceptor{
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
		serverinterceptors.MetricsInterceptor,
		serverinterceptors.LoggingInterceptor,
	}

	resultInterceptors = append(resultInterceptors, interceptors...)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				resultInterceptors...,
			),
		),
	)

	reflection.Register(s)

	return &grpcServer{
		lis:        lis,
		grpcServer: s,
	}
}

func (s *grpcServer) Serve() error {
	return s.grpcServer.Serve(s.lis)
}

func (s *grpcServer) GetServer() *grpc.Server {
	return s.grpcServer
}

func (s *grpcServer) GetListener() net.Listener {
	return s.lis
}
