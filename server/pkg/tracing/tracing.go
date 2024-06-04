package tracing

import (
	"context"
	"os"

	"github.com/uber/jaeger-client-go/config"
)

type CloseFunc func(ctx context.Context) error

func Init(serviceName string) (CloseFunc, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: os.Getenv("JAEGER_HOST"),
		},
	}

	close, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		return close.Close()
	}, nil
}
