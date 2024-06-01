package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
	"github.com/bufbuild/protovalidate-go"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors []grpc.UnaryServerInterceptor
	UnaryInterceptors      []grpc.UnaryServerInterceptor
}

// Deps - server deps
type Deps struct {
	ChannelUsecase usecases.UsecaseInterface
	Log            *log.Logger
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

func NewChannelServer(ctx context.Context, d Deps) (*ChannelServer, error) {
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

	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
