package servers

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

type Grpc struct {
	srv  *grpc.Server
	port string
}

func NewGrpc(srv *grpc.Server, port string) *Grpc {
	return &Grpc{
		srv:  srv,
		port: port,
	}
}

func (s *Grpc) Start() error {
	log.Println("start grpc server", s.port)
	reflection.Register(s.srv)

	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatalf("server: failed to listen: %v", err)
	}

	log.Println("start serve grpc server", lis.Addr())

	go func() {
		if err := s.srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: serve grpc: %v", err)
		}
	}()

	return nil
}

func (s *Grpc) Stop(_ context.Context) error {
	defer log.Println("gateway server is stopped")
	s.srv.GracefulStop()
	return nil
}
