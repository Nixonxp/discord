package servers

import (
	"context"
	"errors"
	"log"
	"net/http"
)

type Gateway struct {
	gatewayHttpServer *http.Server
	port              string
}

func NewGateway(gatewayHttpServer *http.Server, port string) *Gateway {
	return &Gateway{
		gatewayHttpServer: gatewayHttpServer,
		port:              port,
	}
}

func (s *Gateway) Start() error {
	log.Println("start server", s.port)
	var err error
	go func() {
		err = s.gatewayHttpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: serve grpc gateway: %v", err)
		}
	}()

	return nil
}

func (s *Gateway) Stop(ctx context.Context) error {
	defer log.Println("gateway server is stopped")
	return s.gatewayHttpServer.Shutdown(ctx)
}
