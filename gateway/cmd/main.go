package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port           = ":8443"
	certFilePath   = "./cert/server-cert.crt"
	keyFilePath    = "./cert/server-key.key"
	caCertFilePath = "./cert/ca-cert.crt"
)

func CreateServerTLSConfig(caCertFilePath, certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509: %v", err)
	}

	clientCA, err := os.ReadFile(caCertFilePath)
	if err != nil {
		return nil, err
	}

	clientCAs := x509.NewCertPool()
	if !clientCAs.AppendCertsFromPEM(clientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Create tls config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	}

	return tlsConfig, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer resourcesShutdownCtxCancel()

	tlsConfig, err := CreateServerTLSConfig(caCertFilePath, certFilePath, keyFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	config := application.Config{
		GRPCGatewayPort:   port,
		SwaggerUIHTTPPort: ":8081",
		DebugPort:         ":8084",
	}

	log, err := logger.NewLogger(logCfg.NewDefaultConfig())
	if err != nil {
		panic("error init logger")
	}

	closer, err := jaeger_tracing.Init("gateway service")
	if err != nil {
		log.Fatal(ctx, err)
	}
	defer closer(resourcesShutdownCtx)

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
	}

	cb := circuitbreaker.New(
		circuitbreaker.WithClock(clock.New()),
		circuitbreaker.WithFailOnContextCancel(true),
		circuitbreaker.WithFailOnContextDeadline(true),
		circuitbreaker.WithHalfOpenMaxSuccesses(1000),
		circuitbreaker.WithOpenTimeoutBackOff(backoff.NewExponentialBackOff()),
		circuitbreaker.WithOpenTimeout(10*time.Second),
		circuitbreaker.WithCounterResetInterval(10*time.Second),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncFailureRate(10000, 0.1)),
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
		"localhost:8180", grpc.WithIdleTimeout(15*time.Second), // todo return
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
		"localhost:8480", // todo return
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
		"localhost:8580", // todo return
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
		Log:         log,
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
		Handler: middleware.JWTMiddleware(
			middleware.PrometheusMiddleware(
				middleware.TracingWrapper(middleware.AllowCORS(mux)),
			),
		),
		TLSConfig: tlsConfig,
		Addr:      port,
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
