package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
)

// Config - server config
type Config struct {
	GRPCPort string
	HTTPPort string

	ChainUnaryInterceptors []grpc.UnaryServerInterceptor
	UnaryInterceptors      []grpc.UnaryServerInterceptor
}

// Deps - server deps
type Deps struct {
	ChannelUsecase usecases.UsecaseInterface
}

type ChannelServer struct {
	pb.UnimplementedChannelServiceServer
	Deps

	validator *protovalidate.Validator

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	http struct {
		lis  *echo.Echo
		port string
	}
}

func NewChannelServer(ctx context.Context, cfg Config, d Deps) (*ChannelServer, error) {
	srv := &ChannelServer{
		Deps: d,
	}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.AddChannelRequest{},
				&pb.DeleteChannelRequest{},
				&pb.JoinChannelRequest{},
				&pb.LeaveChannelRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}

	// grpc
	{
		grpcServerOptions := unaryInterceptorsToGrpcServerOptions(cfg.UnaryInterceptors...)
		grpcServerOptions = append(grpcServerOptions,
			grpc.ChainUnaryInterceptor(cfg.ChainUnaryInterceptors...),
		)

		grpcServer := grpc.NewServer(grpcServerOptions...)
		pb.RegisterChannelServiceServer(grpcServer, srv)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			return nil, fmt.Errorf("server: failed to listen: %v", err)
		}

		srv.grpc.server = grpcServer
		srv.grpc.lis = lis
	}

	{
		e := echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		e.GET("/health", func(c echo.Context) error {
			status := http.StatusOK
			statusMessage := "OK"

			if !isServiceOk(10) {
				status = http.StatusInternalServerError
				statusMessage = "Error"
			}

			return c.JSON(status, struct{ Status string }{Status: statusMessage})
		})

		e.GET("/ready", func(c echo.Context) error {
			status := http.StatusOK
			statusMessage := "OK"

			if !isServiceOk(5) {
				status = http.StatusInternalServerError
				statusMessage = "Error"
			}

			return c.JSON(status, struct{ Status string }{Status: statusMessage})
		})

		srv.http.lis = e
		srv.http.port = cfg.HTTPPort
	}

	return srv, nil
}

// Run - serve
func (s *ChannelServer) Run(ctx context.Context) error {
	group := errgroup.Group{}

	group.Go(func() error {
		log.Println("start serve", s.grpc.lis.Addr())
		if err := s.grpc.server.Serve(s.grpc.lis); err != nil {
			return fmt.Errorf("server: serve grpc: %v", err)
		}
		return nil
	})

	group.Go(func() error {
		log.Println("start http server", s.http.port)
		err := s.http.lis.Start(s.http.port)
		if err != nil {
			return fmt.Errorf("server: serve http: %v", err)
		}

		return nil
	})

	return group.Wait()
}

func unaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}

// isServiceOk в зависимости от входящего значения вернет false, например
// передано 5, тогда (100 / 5 = 20) 20% вероятностью вернется false, для теста сервиса
func isServiceOk(probability int) bool {
	randNumber := rand.IntN(probability-1) + 1

	if randNumber == 1 {
		return false
	}

	return true
}
