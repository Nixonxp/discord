package main

import (
	"context"
	repository "github.com/Nixonxp/discord/auth/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/auth/internal/app/server"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/auth/internal/middleware/errors"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	userInMemoryRepo := repository.NewInMemoryUserRepository()
	authUsecase := usecases.NewAuthUsecase(usecases.Deps{
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

	srv, err := server.NewAuthServer(ctx, config, server.Deps{
		AuthUsecase: authUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
