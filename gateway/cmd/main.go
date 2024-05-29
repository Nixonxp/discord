package main

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/gateway/internal/app/server"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	"github.com/Nixonxp/discord/gateway/internal/middleware"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	"github.com/Nixonxp/discord/gateway/pkg/application"
	logCfg "github.com/Nixonxp/discord/gateway/pkg/logger"
	logger "github.com/Nixonxp/discord/gateway/pkg/logger"
	jaeger_tracing "github.com/Nixonxp/discord/gateway/pkg/tracing"
	"github.com/benbjohnson/clock"
	"github.com/cenkalti/backoff/v3"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/mercari/go-circuitbreaker"
	"github.com/opentracing/opentracing-go"
	ratelimit "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
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

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer resourcesShutdownCtxCancel()

	config := application.Config{
		GRPCGatewayPort:   ":8080",
		SwaggerUIHTTPPort: ":8081",
	}

	log, err := logger.NewLogger(logCfg.NewDefaultConfig())
	if err != nil {
		panic("error init logger")
	}

	if err := jaeger_tracing.Init("gateway service"); err != nil {
		log.Fatal(ctx, err)
	}

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
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

	srvCfg := server.Config{
		UnaryClientInterceptors: []grpc.UnaryClientInterceptor{
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpc_opentracing.OpenTracingClientInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware.UnaryCircuitBreakerClientInterceptor(
				cb,
				func(ctx context.Context, method string, req interface{}) {
					fmt.Printf("[Client] Circuit breaker is open.\n")
				},
			),
			ratelimit.UnaryClientInterceptor(ratelimit.NewLimiter(10000)),
		},
	}

	// grpc connection to auth service
	authConn, err := grpc.DialContext(resourcesShutdownCtx,
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
	userConn, err := grpc.DialContext(resourcesShutdownCtx,
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
	serverConn, err := grpc.DialContext(resourcesShutdownCtx,
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
	channelConn, err := grpc.DialContext(resourcesShutdownCtx,
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
	chatConn, err := grpc.DialContext(resourcesShutdownCtx,
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

	srv, err := server.NewDiscordGatewayServiceServer(resourcesShutdownCtx, server.Deps{
		DiscordGatewayService: discordGatewayService,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	mux := runtime.NewServeMux()

	if err := pb.RegisterGatewayServiceHandlerServer(resourcesShutdownCtx, mux, srv); err != nil {
		log.Fatalf("server: failed to register handler: %v", err)
	}

	httpServer := &http.Server{
		Handler: middleware.TracingWrapper(middleware.AllowCORS(mux)),
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	app.AddGatewayServer(httpServer)

	if err = app.Run(ctx, nil); err != nil {
		log.Fatalf("run: %v", err)
	}

	log.Print("servers is stopped")
	resourcesShutdownCtxCancel()
	log.Print("wait shutdown grpc clients")
	time.Sleep(time.Second * 5)

	defer log.Print("app is stopped")
}
