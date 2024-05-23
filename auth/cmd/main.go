package main

import (
	"context"
	repository "github.com/Nixonxp/discord/auth/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/auth/internal/app/server"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/auth/internal/middleware/errors"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	"github.com/Nixonxp/discord/auth/pkg/application"
	pkg_middleware "github.com/Nixonxp/discord/auth/pkg/middleware"
	"github.com/Nixonxp/discord/auth/pkg/rate_limiter"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := application.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	userInMemoryRepo := repository.NewInMemoryUserRepository()
	authUsecase := usecases.NewAuthUsecase(usecases.Deps{
		UserRepo: userInMemoryRepo,
	})

	// limiter per method
	methodLimiter := pkg_middleware.NewMethodRateLimiterInterceptor(
		pkg_middleware.NewMethodLimiterInfo("register", 100),
		pkg_middleware.NewMethodLimiterInfo("login", 500),
	)

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
			ratelimit.UnaryServerInterceptor(globalLimiter),
			methodLimiter.GetInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		},
	}

	srv, err := server.NewAuthServer(ctx, server.Deps{
		AuthUsecase: authUsecase,
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
}
