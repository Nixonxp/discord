package main

import (
	"context"
	repository "github.com/Nixonxp/discord/auth/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/auth/internal/app/server"
	"github.com/Nixonxp/discord/auth/internal/app/services/user"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/auth/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/auth/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/auth/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	"github.com/Nixonxp/discord/auth/pkg/application"
	logCfg "github.com/Nixonxp/discord/auth/pkg/logger"
	logger "github.com/Nixonxp/discord/auth/pkg/logger"
	method_rate_limiter "github.com/Nixonxp/discord/auth/pkg/method_rate_limiter"
	"github.com/Nixonxp/discord/auth/pkg/rate_limiter"
	jaeger_tracing "github.com/Nixonxp/discord/auth/pkg/tracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	ratelimitCustom "github.com/tommy-sho/rate-limiter-grpc-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	log, err := logger.NewLogger(logCfg.NewDefaultConfig())
	if err != nil {
		panic("error init logger")
	}

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)

	closer, err := jaeger_tracing.Init("auth service")
	if err != nil {
		log.Fatal(ctx, err)
	}
	defer closer(resourcesShutdownCtx)

	config := application.Config{
		GRPCPort:  ":8080",
		HTTPPort:  ":8081",
		DebugPort: ":8084",
	}

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
	}

	serverConfig := server.Config{
		UserServiceUrl: "user:8080",
		UnaryClientInterceptors: []grpc.UnaryClientInterceptor{
			grpcretry.UnaryClientInterceptor(retryOpts...),
			ratelimitCustom.UnaryClientInterceptor(ratelimitCustom.NewLimiter(1000)),
			grpc_opentracing.OpenTracingClientInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
		},
	}

	userServiceClient, err := user.NewClient(serverConfig, log)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     "586587681748-68glr3bqvfu64q5jks769ba78pjhp309.apps.googleusercontent.com", // todo to env file
		ClientSecret: "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",
		RedirectURL:  "https://localhost:8443/api/v1/oauth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	userInMemoryRepo := repository.NewInMemoryUserRepository()
	authUsecase := usecases.NewAuthUsecase(usecases.Deps{
		UserRepo:    userInMemoryRepo,
		UserService: userServiceClient,
		Log:         log,
		Oauth2Cgf:   oauthConfig,
	})

	// limiter per method
	methodLimiter := method_rate_limiter.NewMethodRateLimiterInterceptor(
		method_rate_limiter.NewMethodLimiterInfo("register", 10000),
		method_rate_limiter.NewMethodLimiterInfo("login", 50000),
	)

	globalLimiter := rate_limiter.NewRateLimiter(100000)
	grpcConfig := server.Config{
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

	srv, err := server.NewAuthServer(resourcesShutdownCtx, server.Deps{
		AuthUsecase: authUsecase,
		Log:         log,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterAuthServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}

	log.Print("servers is stopped")

	log.Print("wait shutdown resources")
	resourcesShutdownCtxCancel()
	time.Sleep(time.Second * 5)

	defer log.Print("app is stopped")
}
