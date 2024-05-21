package application

import (
	"context"
)

type Server interface {
	Run(ctx context.Context) error
}

type App struct {
	server Server
}

func NewApp(server Server) (App, error) {
	return App{
		server: server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.server.Run(ctx); err != nil {
		return err
	}

	return nil
}
