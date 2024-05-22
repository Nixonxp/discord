package main

import (
	"context"
	"github.com/Nixonxp/discord/gateway/internal/app/server"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	"github.com/Nixonxp/discord/gateway/internal/middleware"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	"github.com/Nixonxp/discord/gateway/pkg/application"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

var digContainer *dig.Container

func ProvideConfig() *application.Config {
	config := application.Config{
		GRPCGatewayPort:   ":8080",
		SwaggerUIHTTPPort: ":8081",
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

func ProvideServer(service *services.DiscordGatewayService) *http.Server {
	srv, err := server.NewDiscordGatewayServiceServer(context.Background(), server.Deps{
		DiscordGatewayService: service,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterGatewayServiceHandlerServer(context.Background(), mux, srv); err != nil {
		log.Fatalf("server: failed to register handler: %v", err)
	}

	httpServer := &http.Server{Handler: middleware.AllowCORS(mux)}

	return httpServer
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
		cfg *application.Config,
		srv *http.Server,
	) {
		app, err := application.NewApp(cfg)
		if err != nil {
			log.Fatalf("failed to create app: %v", err)
		}

		app.AddGatewayServer(srv)

		if err = app.Run(ctx, nil); err != nil {
			log.Fatalf("run: %v", err)
		}
	}

	if err := digContainer.Invoke(applicationRun); err != nil {
		log.Fatalf("run: %v", err)
	}
}
