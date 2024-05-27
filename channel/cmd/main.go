package main

import (
	"context"
	repository "github.com/Nixonxp/discord/channel/internal/app/repository/channel_storage"
	"github.com/Nixonxp/discord/channel/internal/app/server"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/channel/internal/middleware/errors"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	"github.com/Nixonxp/discord/channel/pkg/application"
	"github.com/Nixonxp/discord/channel/pkg/rate_limiter"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"log"
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

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer resourcesShutdownCtxCancel()

	config := application.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	channelInMemoryRepo := repository.NewInMemoryChannelRepository()
	authUsecase := usecases.NewChannelUsecase(usecases.Deps{
		ChannelRepo: channelInMemoryRepo,
	})

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
		},
	}

	srv, err := server.NewChannelServer(resourcesShutdownCtx, server.Deps{
		ChannelUsecase: authUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterChannelServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}

	log.Print("servers is stopped")
	resourcesShutdownCtxCancel()
	log.Print("wait shutdown resources")
	time.Sleep(time.Second * 5)

	defer log.Print("app is stopped")
}
