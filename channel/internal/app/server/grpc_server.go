package server

import (
	"context"
	"fmt"
	repository "github.com/Nixonxp/discord/channel/internal/app/repository/channel_storage"
	sub_repository "github.com/Nixonxp/discord/channel/internal/app/repository/subscribe_storage"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/channel/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/channel/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/channel/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
	"github.com/Nixonxp/discord/channel/pkg/rate_limiter"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
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

	grpcServer *grpc.Server
}

func NewChannelServer(ctx context.Context, s *MainServer) (*ChannelServer, error) {
	srv := &ChannelServer{}

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

	mainCollection := s.mongo.GetInstance()
	subscribeCollection, err := mainCollection.NewCollection(s.cfg.Application.ChannelSubscribeCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect mongo: %v", err)
	}

	channelRepo := repository.NewMongoChannelRepository(mainCollection, s.logger.GetInstance())
	subscribeRepo := sub_repository.NewMongoSubscribeRepository(subscribeCollection, s.logger.GetInstance())
	authUsecase := usecases.NewChannelUsecase(usecases.Deps{
		ChannelRepo:   channelRepo,
		SubscribeRepo: subscribeRepo,
	})

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true),
			middleware_metrics.MetricsUnaryInterceptor(),
		},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv.ChannelUsecase = authUsecase
	srv.Log = s.logger.GetInstance()

	grpcServerOptions := UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterChannelServiceServer(grpcServer, srv)

	srv.grpcServer = grpcServer

	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
