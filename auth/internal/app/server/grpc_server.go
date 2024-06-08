package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/auth/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/auth/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/auth/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	"github.com/Nixonxp/discord/auth/pkg/method_rate_limiter"
	"github.com/Nixonxp/discord/auth/pkg/rate_limiter"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors  []grpc.UnaryServerInterceptor
	UnaryInterceptors       []grpc.UnaryServerInterceptor
	UnaryClientInterceptors []grpc.UnaryClientInterceptor
}

// Deps - server deps
type Deps struct {
	AuthUsecase usecases.UsecaseInterface
	Log         *log.Logger
}

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	Deps

	validator  *protovalidate.Validator
	grpcServer *grpc.Server

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	http struct {
		lis  *echo.Echo
		port string
	}
}

func NewAuthServer(ctx context.Context, s *MainServer) (*AuthServer, error) {
	srv := &AuthServer{}
	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.RegisterRequest{},
				&pb.LoginRequest{},
				&pb.OauthLoginRequest{},
				&pb.OauthLoginCallbackRequest{},
				&pb.RefreshRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}

	methodLimiter := method_rate_limiter.NewMethodRateLimiterInterceptor(
		method_rate_limiter.NewMethodLimiterInfo("register", 10000),
		method_rate_limiter.NewMethodLimiterInfo("login", 50000),
	)

	globalLimiter := rate_limiter.NewRateLimiter(100000)
	grpcConfig := Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			ratelimit.UnaryServerInterceptor(globalLimiter),
			methodLimiter.GetInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true),
			middleware_metrics.MetricsUnaryInterceptor(),
		},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	oauthConfig := oauth2.Config{
		ClientID:     s.cfg.Application.OAuthClientID,
		ClientSecret: s.cfg.Application.OAuthClientSecret,
		RedirectURL:  s.cfg.Application.OAuthRedirectUrl,
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	authUsecase := usecases.NewAuthUsecase(usecases.Deps{
		UserService: s.userSvcClient.GetInstance(),
		Log:         s.logger.GetInstance(),
		Oauth2Cgf:   oauthConfig,
	})

	grpcServerOptions := UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	srv.AuthUsecase = authUsecase
	srv.Log = s.logger.GetInstance()

	srv.grpcServer = grpc.NewServer(grpcServerOptions...)
	pb.RegisterAuthServiceServer(srv.grpcServer, srv)

	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
