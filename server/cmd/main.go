package main

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/repository/server_storage"
	"github.com/Nixonxp/discord/server/internal/app/server"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/server/internal/middleware/errors"
	pb "github.com/Nixonxp/discord/server/pkg/api/v1"
	"github.com/Nixonxp/discord/server/pkg/application"
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

	serverInMemoryRepo := repository.NewInMemoryServerRepository()
	serverUsecase := usecases.NewServerUsecase(usecases.Deps{
		ServerRepo: serverInMemoryRepo,
	})

	// delivery
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewServerServer(ctx, server.Deps{
		ServerUsecase: serverUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterServerServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}
}
