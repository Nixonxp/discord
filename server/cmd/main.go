package main

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/repository/server_storage"
	"github.com/Nixonxp/discord/server/internal/app/server"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/server/internal/middleware/errors"
	"github.com/Nixonxp/discord/server/pkg/application"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	serverInMemoryRepo := repository.NewInMemoryServerRepository()
	serverUsecase := usecases.NewServerUsecase(usecases.Deps{
		ServerRepo: serverInMemoryRepo,
	})

	// delivery
	config := server.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewServerServer(ctx, config, server.Deps{
		ServerUsecase: serverUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	app, err := application.NewApp(srv)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err = app.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
