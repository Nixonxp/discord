package server

import (
	"context"
	"fmt"
	config "github.com/Nixonxp/discord/gateway/configs"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	gateway_serv "github.com/Nixonxp/discord/gateway/internal/app/services/gateway"
	"github.com/Nixonxp/discord/gateway/internal/middleware"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	"github.com/Nixonxp/discord/gateway/pkg/servers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"sync"
	"time"
)

const port = ":8443"

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
	gateway gateway_serv.Gateway
	servers []Server
	cfg     *config.ApplicationConfig
}

func (s *MainServer) AddServer(srv Server) {
	s.servers = append(s.servers, srv)
}

func (s *MainServer) AppStart(ctx context.Context, cfg *config.Config) error {
	srv, err := NewDiscordGatewayServiceServer(ctx, Deps{
		DiscordGatewayService: s.gateway.GetInstance(),
	})
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	mux := runtime.NewServeMux()

	if err := pb.RegisterGatewayServiceHandlerServer(ctx, mux, srv); err != nil {
		return fmt.Errorf("server: failed to register handler: %v", err)
	}

	httpServer := &http.Server{
		Handler: middleware.JWTMiddleware(
			middleware.PrometheusMiddleware(
				middleware.TracingWrapper(middleware.AllowCORS(mux)),
			),
		),
		Addr: cfg.Application.Port,
	}

	gatewayServer := servers.NewGateway(httpServer, cfg.Application.Port)
	s.AddServer(gatewayServer)

	swaggerui := servers.NewSwagger(cfg.Application.SwaggerPort)
	s.AddServer(swaggerui)

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
