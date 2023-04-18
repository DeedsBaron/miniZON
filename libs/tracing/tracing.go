package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"route256/libs/logger"
)

func Init(serviceName string) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	cfg2, err := cfg.FromEnv()
	if err != nil {
		logger.Fatalf("Cannot init tracing %v", err)
	}
	tracer, _, err := cfg2.NewTracer()
	if err != nil {
		logger.Fatalf("Cannot init tracing %v", err)
	}
	opentracing.SetGlobalTracer(tracer)
}
