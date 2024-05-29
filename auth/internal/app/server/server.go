package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	"github.com/bufbuild/protovalidate-go"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors  []grpc.UnaryServerInterceptor
	UnaryInterceptors       []grpc.UnaryServerInterceptor
	UnaryClientInterceptors []grpc.UnaryClientInterceptor

	UserServiceUrl string
}

// Deps - server deps
type Deps struct {
	AuthUsecase usecases.UsecaseInterface
	Log         *log.Logger
}

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
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

func NewAuthServer(ctx context.Context, d Deps) (*AuthServer, error) {
	srv := &AuthServer{
		Deps: d,
	}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.RegisterRequest{},
				&pb.LoginRequest{},
				&pb.OauthLoginRequest{},
				&pb.OauthLoginCallbackRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}
	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
