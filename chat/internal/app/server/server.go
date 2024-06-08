package server

import (
	"context"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/Nixonxp/discord/chat/internal/app/services"
	"github.com/Nixonxp/discord/chat/pkg/servers"
	"sync"
	"time"
)

type Server interface {
	// Start starts server in separate goroutine, method is non-blocking.
	Start() error
	// Stop stops currently running server and blocks until server is stopped.
	Stop(ctx context.Context) error
}

var terminationTimeout = time.Second * 10

type MainServer struct {
	tracer  services.Tracing
	logger  services.Logger
	mongo   services.Mongo
	kafka   services.Kafka
	servers []Server
	cfg     *config.Config
}

func (s *MainServer) AddServer(srv Server) {
	s.servers = append(s.servers, srv)
}

func (s *MainServer) AppStart(ctx context.Context, cfg *config.Config) error {
	s.cfg = cfg

	srv, err := NewChatServer(ctx, s)
	if err != nil {
		s.logger.GetInstance().Fatalf("failed to create server: %v", err)
	}

	server := servers.NewGrpc(srv.grpcServer, s.cfg.Application.Port)
	s.AddServer(server)

	metrics := servers.NewMetrics(cfg.Application.MetricsPort)
	s.AddServer(metrics)

	// start servers
	for _, srvItem := range s.servers {
		if err := srvItem.Start(); err != nil {
			s.logger.GetInstance().Info("failed to start app server")
			return err
		}
	}

	// listen for context cancel
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), terminationTimeout)
	defer cancel()
	// stop servers in parallel
	wg := new(sync.WaitGroup)
	for _, srvItem := range s.servers {
		wg.Add(1)
		go func(srvItem Server) {
			defer wg.Done()
			if err := srvItem.Stop(shutdownCtx); err != nil {
				s.logger.GetInstance().WithError(err).Info("failed to stop app server")
			}
		}(srvItem)
	}
	wg.Wait()

	return nil
}
