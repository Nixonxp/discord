package services

import (
	"context"
	config "github.com/Nixonxp/discord/gateway/configs"
	logCfg "github.com/Nixonxp/discord/gateway/pkg/logger"
	logger "github.com/Nixonxp/discord/gateway/pkg/logger"
	"sync"
)

type Logger struct {
}

var instance *logger.Logger
var once sync.Once

func (l *Logger) Init(_ context.Context, _ *config.Config) error {
	var err error
	once.Do(func() {
		instance, err = logger.NewLogger(logCfg.NewDefaultConfig())
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *Logger) Ident() string {
	return "logger"
}

func (l *Logger) GetInstance() *logger.Logger {
	return instance
}

func (l *Logger) Close(_ context.Context) error {
	return nil
}
