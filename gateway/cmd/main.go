package main

import (
	"context"
	"github.com/Nixonxp/discord/gateway/internal/app/server"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	"github.com/Nixonxp/discord/gateway/pkg/application"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var digContainer *dig.Container

func ProvideConfig() *server.Config {
	config := server.Config{
		GRPCGatewayPort:        ":8080",
		HTTPSwaggerUIPort:      ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{},
	}

	return &config
}

func ProvideGatewayService() *services.DiscordGatewayService {
	// grpc connection to auth service
	authConn, err := grpc.DialContext(context.Background(), "auth:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to user service
	userConn, err := grpc.DialContext(context.Background(), "user:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to server service
	serverConn, err := grpc.DialContext(context.Background(), "server:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to channel service
	channelConn, err := grpc.DialContext(context.Background(), "channel:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to chat service
	chatConn, err := grpc.DialContext(context.Background(), "chat:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	discordGatewayService := services.NewDiscordGatewayService(services.Deps{
		AuthConn:    authConn,
		UserConn:    userConn,
		ServerConn:  serverConn,
		ChannelConn: channelConn,
		ChatConn:    chatConn,
	})

	return discordGatewayService
}

func ProvideServer(config *server.Config, service *services.DiscordGatewayService) *server.DiscordGatewayServiceServer {
	srv, err := server.NewDiscordGatewayServiceServer(context.Background(), *config, server.Deps{
		DiscordGatewayService: service,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	return srv
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	digContainer = dig.New()
	digContainer.Provide(ProvideConfig)
	digContainer.Provide(ProvideGatewayService)
	digContainer.Provide(ProvideServer)

	applicationRun := func(
		srv *server.DiscordGatewayServiceServer,
	) {
		if err := srv.Run(ctx); err != nil {
			log.Fatalf("run: %v", err)
		}

		app, err := application.NewApp(srv)
		if err != nil {
			log.Fatalf("failed to create app: %v", err)
		}

		if err = app.Run(ctx); err != nil {
			log.Fatalf("run: %v", err)
		}
	}

	if err := digContainer.Invoke(applicationRun); err != nil {
		log.Fatalf("run: %v", err)
	}
}
