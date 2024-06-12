package services

import (
	"context"
	config "github.com/Nixonxp/discord/chat/configs"
	jaeger_tracing "github.com/Nixonxp/discord/chat/pkg/tracing"
)

type Tracing struct {
	closer jaeger_tracing.CloseFunc
}

func (t *Tracing) Init(_ context.Context, _ *config.Config) error {
	var err error
	closer, err := jaeger_tracing.Init("chat service")
	if err != nil {
		return err
	}

	t.closer = closer
	return nil
}

func (t *Tracing) Ident() string {
	return "tracing"
}

func (t *Tracing) Close(ctx context.Context) error {
	err := t.closer(ctx)
	if err != nil {
		return err
	}
	return nil
}
