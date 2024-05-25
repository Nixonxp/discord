package main

import (
	"context"
	repository "github.com/Nixonxp/discord/user/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/user/internal/app/server"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/user/internal/middleware/errors"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	"github.com/Nixonxp/discord/user/pkg/application"
	"github.com/Nixonxp/discord/user/pkg/postgres"
	"github.com/Nixonxp/discord/user/pkg/rate_limiter"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"log"
	"time"
)

const DSN = "user=admin password=password123 host=postgres port=5432 dbname=discord sslmode=require pool_max_conns=10"

//const DSN = "user=admin password=password123 host=localhost port=5432 dbname=discord sslmode=require pool_max_conns=10" // todo delete

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

	// repository
	pool, err := postgres.NewConnectionPool(ctx, DSN,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserPostgresqlRepository(pool)

	userUsecase := usecases.NewUserUsecase(usecases.Deps{
		UserRepo: userRepo,
	})

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
		},
	}

	srv, err := server.NewUserServer(ctx, server.Deps{
		UserUsecase: userUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterUserServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}
}
