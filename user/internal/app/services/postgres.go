package services

import (
	"context"
	config "github.com/Nixonxp/discord/user/configs"
	"github.com/Nixonxp/discord/user/pkg/postgres"
	"github.com/Nixonxp/discord/user/pkg/transaction_manager"
	"time"
)

type Postgres struct {
	pool      *postgres.Connection
	txManager *transaction_manager.TransactionManager
}

func (p *Postgres) Init(ctx context.Context, cfg *config.Config) error {
	var err error
	// repository
	p.pool, err = postgres.NewConnectionPool(ctx, cfg.Application.PostgresDSN,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		return err
	}

	p.txManager = transaction_manager.New(p.pool)
	return nil
}

func (p *Postgres) GetInstance() *transaction_manager.TransactionManager {
	return p.txManager
}

func (p *Postgres) Ident() string {
	return "postgres"
}

func (p *Postgres) Close(_ context.Context) error {
	err := p.pool.Close()
	if err != nil {
		return err
	}
	return nil
}
