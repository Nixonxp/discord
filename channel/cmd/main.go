package main

import (
	"context"
	repository "github.com/Nixonxp/discord/channel/internal/app/repository/channel_storage"
	sub_repository "github.com/Nixonxp/discord/channel/internal/app/repository/subscribe_storage"
	"github.com/Nixonxp/discord/channel/internal/app/server"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/channel/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/channel/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/channel/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	"github.com/Nixonxp/discord/channel/pkg/application"
	logCfg "github.com/Nixonxp/discord/channel/pkg/logger"
	logger "github.com/Nixonxp/discord/channel/pkg/logger"
	mongoCollection "github.com/Nixonxp/discord/channel/pkg/mongo"
	"github.com/Nixonxp/discord/channel/pkg/rate_limiter"
	jaeger_tracing "github.com/Nixonxp/discord/channel/pkg/tracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"os"
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
		GRPCPort:  ":8080",
		HTTPPort:  ":8081",
		DebugPort: ":8084",
	}

	log, err := logger.NewLogger(logCfg.NewDefaultConfig())
	if err != nil {
		panic("error init logger")
	}

	closer, err := jaeger_tracing.Init("channel service")
	if err != nil {
		log.Fatal(ctx, err)
	}
	defer closer(resourcesShutdownCtx)

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	collection, err := mongoCollection.NewCollection(resourcesShutdownCtx,
		"channels", // todo from config
		&mongoCollection.Config{
			MongoHost:     os.Getenv("MONGO_HOST"),
			MongoDb:       os.Getenv("MONGO_DB"),
			MongoPort:     os.Getenv("MONGO_PORT"),
			MongoUser:     os.Getenv("MONGO_USER"),
			MongoPassword: os.Getenv("MONGO_PASSWORD"),
		},
	)

	subscribeCollection, err := collection.NewCollection("channel_subscribe") // todo to config
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	channelRepo := repository.NewMongoChannelRepository(collection, log)
	subscribeRepo := sub_repository.NewMongoSubscribeRepository(subscribeCollection, log)
	authUsecase := usecases.NewChannelUsecase(usecases.Deps{
		ChannelRepo:   channelRepo,
		SubscribeRepo: subscribeRepo,
	})

	globalLimiter := rate_limiter.NewRateLimiter(1000)
	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true),
			middleware_metrics.MetricsUnaryInterceptor(),
		},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	srv, err := server.NewChannelServer(resourcesShutdownCtx, server.Deps{
		ChannelUsecase: authUsecase,
		Log:            log,
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
