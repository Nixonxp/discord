package main

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/gateway/internal/app/server"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	"github.com/Nixonxp/discord/gateway/internal/middleware"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	"github.com/Nixonxp/discord/gateway/pkg/application"
	"github.com/benbjohnson/clock"
	"github.com/cenkalti/backoff/v3"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mercari/go-circuitbreaker"
	ratelimit "github.com/tommy-sho/rate-limiter-grpc-go"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"time"
)

var digContainer *dig.Container

func ProvideConfig() (*application.Config, *server.Config) {
	config := application.Config{
		GRPCGatewayPort:   ":8080",
		SwaggerUIHTTPPort: ":8081",
	}

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 5),
	}

	cb := circuitbreaker.New(
		circuitbreaker.WithClock(clock.New()),
		circuitbreaker.WithFailOnContextCancel(true),
		circuitbreaker.WithFailOnContextDeadline(true),
		circuitbreaker.WithHalfOpenMaxSuccesses(100),
		circuitbreaker.WithOpenTimeoutBackOff(backoff.NewExponentialBackOff()),
		circuitbreaker.WithOpenTimeout(10*time.Second),
		circuitbreaker.WithCounterResetInterval(10*time.Second),
		// we also have NewTripFuncThreshold and NewTripFuncConsecutiveFailures
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncFailureRate(100, 0.1)),
		circuitbreaker.WithOnStateChangeHookFn(func(from, to circuitbreaker.State) {
			log.Printf("state changed from %s to %s\n", from, to)
		}),
	)

	serverConfig := server.Config{
		UnaryClientInterceptors: []grpc.UnaryClientInterceptor{
			grpcretry.UnaryClientInterceptor(retryOpts...),
			middleware.UnaryCircuitBreakerClientInterceptor(
				cb,
				func(ctx context.Context, method string, req interface{}) {
					fmt.Printf("[Client] Circuit breaker is open.\n")
				},
			),
			ratelimit.UnaryClientInterceptor(ratelimit.NewLimiter(10000)),
		},
	}

	return &config, &serverConfig
}

func ProvideGatewayService(srvCfg *server.Config) *services.DiscordGatewayService {
	// grpc connection to auth service
	authConn, err := grpc.DialContext(context.Background(),
		"auth:8080", grpc.WithIdleTimeout(15*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to user service
	userConn, err := grpc.DialContext(context.Background(),
		"user:8080", grpc.WithIdleTimeout(15*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to server service
	serverConn, err := grpc.DialContext(context.Background(),
		"server:8080",
		grpc.WithIdleTimeout(15*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to channel service
	channelConn, err := grpc.DialContext(context.Background(),
		"channel:8080", grpc.WithIdleTimeout(15*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	// grpc connection to chat service
	chatConn, err := grpc.DialContext(context.Background(),
		"chat:8080",
		grpc.WithIdleTimeout(15*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
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
