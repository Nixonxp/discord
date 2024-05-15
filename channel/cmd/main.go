package main

import (
	"context"
	repository "github.com/Nixonxp/discord/channel/internal/app/repository/channel_storage"
	"github.com/Nixonxp/discord/channel/internal/app/server"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/channel/internal/middleware/errors"
	"google.golang.org/grpc"
	"log"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	channelInMemoryRepo := repository.NewInMemoryChannelRepository()
	authUsecase := usecases.NewChannelUsecase(usecases.Deps{
		ChannelRepo: channelInMemoryRepo,
	})

	// delivery
	config := server.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewChannelServer(ctx, config, server.Deps{
		ChannelUsecase: authUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
