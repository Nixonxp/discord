package main

import (
	"context"
	"github.com/Nixonxp/discord/gateway/internal/app/server"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	// delivery
	config := server.Config{
		GRPCGatewayPort:        ":8080",
		HTTPSwaggerUIPort:      ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{},
	}

	srv, err := server.NewDiscordGatewayServiceServer(ctx, config, server.Deps{
		DiscordGatewayService: discordGatewayService,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}

}
