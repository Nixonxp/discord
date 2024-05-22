package main

import (
	"context"
	repository "github.com/Nixonxp/discord/user/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/user/internal/app/server"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/user/internal/middleware/errors"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	"github.com/Nixonxp/discord/user/pkg/application"
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
	userUsecase := usecases.NewUserUsecase(usecases.Deps{
		UserRepo: userInMemoryRepo,
	})

	// delivery
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
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
