package main

import (
	"context"
	repository "github.com/Nixonxp/discord/user/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/user/internal/app/server"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/user/internal/middleware/errors"
	"github.com/Nixonxp/discord/user/pkg/application"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userInMemoryRepo := repository.NewInMemoryUserRepository()
	userUsecase := usecases.NewUserUsecase(usecases.Deps{
		UserRepo: userInMemoryRepo,
	})

	// delivery
	config := server.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewUserServer(ctx, config, server.Deps{
		UserUsecase: userUsecase,
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
