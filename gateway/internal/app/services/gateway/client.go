package services

import (
	"context"
	"fmt"
	config "github.com/Nixonxp/discord/gateway/configs"
	"github.com/Nixonxp/discord/gateway/internal/app/services"
	"github.com/Nixonxp/discord/gateway/internal/middleware"
	"github.com/benbjohnson/clock"
	"github.com/cenkalti/backoff/v3"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/mercari/go-circuitbreaker"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	ratelimit "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors  []grpc.UnaryServerInterceptor
	UnaryInterceptors       []grpc.UnaryServerInterceptor
	UnaryClientInterceptors []grpc.UnaryClientInterceptor
}

type Gateway struct {
	service *DiscordGatewayService
}

func (g *Gateway) Init(ctx context.Context, cfg *config.Config) error {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
	}

	cb := circuitbreaker.New(
		circuitbreaker.WithClock(clock.New()),
		circuitbreaker.WithFailOnContextCancel(true),
		circuitbreaker.WithFailOnContextDeadline(true),
		circuitbreaker.WithHalfOpenMaxSuccesses(10000),
		circuitbreaker.WithOpenTimeoutBackOff(backoff.NewExponentialBackOff()),
		circuitbreaker.WithOpenTimeout(10*time.Second),
		circuitbreaker.WithCounterResetInterval(10*time.Second),
		circuitbreaker.WithTripFunc(circuitbreaker.NewTripFuncFailureRate(10000, 0.1)),
		circuitbreaker.WithOnStateChangeHookFn(func(from, to circuitbreaker.State) {
			log.Printf("state changed from %s to %s\n", from, to)
		}),
	)

	srvCfg := Config{
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
	authConn, err := grpc.DialContext(ctx,
		cfg.Application.AuthServiceHost, grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	// grpc connection to user service
	userConn, err := grpc.DialContext(ctx,
		cfg.Application.UserServiceHost, grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	// grpc connection to server service
	serverConn, err := grpc.DialContext(ctx,
		cfg.Application.ServerServiceHost,
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	// grpc connection to channel service
	channelConn, err := grpc.DialContext(ctx,
		cfg.Application.ChannelServiceHost, grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	// grpc connection to chat service
	chatConn, err := grpc.DialContext(ctx,
		cfg.Application.ChatServiceHost,
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			srvCfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	logger := services.Logger{}
	err = logger.Init(ctx, cfg)
	if err != nil {
		return err
	}

	g.service = NewDiscordGatewayService(Deps{
		AuthConn:    authConn,
		UserConn:    userConn,
		ServerConn:  serverConn,
		ChannelConn: channelConn,
		ChatConn:    chatConn,
		Log:         logger.GetInstance(),
	})

	return nil
}

func (g *Gateway) Ident() string {
	return "gateway service"
}

func (g *Gateway) GetInstance() *DiscordGatewayService {
	return g.service
}

func (g *Gateway) Close(_ context.Context) error {
	return nil
}
