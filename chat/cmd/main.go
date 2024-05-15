package main

import (
	"context"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	"github.com/Nixonxp/discord/chat/internal/app/server"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/chat/internal/middleware/errors"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chatInMemoryRepo := repository.NewInMemoryChatRepository()
	chatUsecase := usecases.NewChatUsecase(usecases.Deps{
		ChatRepo: chatInMemoryRepo,
	})

	// delivery
	config := server.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewChatServer(ctx, config, server.Deps{
		ChatUsecase: chatUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
