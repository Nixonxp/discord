package main

import (
	"context"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	"github.com/Nixonxp/discord/chat/internal/app/server"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/chat/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/chat/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/chat/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/chat/pkg/api/v1"
	"github.com/Nixonxp/discord/chat/pkg/application"
	logCfg "github.com/Nixonxp/discord/chat/pkg/logger"
	logger "github.com/Nixonxp/discord/chat/pkg/logger"
	mongoCollection "github.com/Nixonxp/discord/chat/pkg/mongo"
	"github.com/Nixonxp/discord/chat/pkg/rate_limiter"
	jaeger_tracing "github.com/Nixonxp/discord/chat/pkg/tracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
	"time"
)

const MongoCollectionMessages = "messages"

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

	closer, err := jaeger_tracing.Init("chat service")
	if err != nil {
		log.Fatal(ctx, err)
	}
	defer closer(resourcesShutdownCtx)

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	collection, err := mongoCollection.NewCollection(resourcesShutdownCtx,
		MongoCollectionMessages,
		&mongoCollection.Config{
			MongoHost:     "mongodb",
			MongoDb:       "Chat",
			MongoPort:     "27017",
			MongoUser:     "root",
			MongoPassword: "example",
		},
	)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	defer collection.DisconnectClient()

	chatMongoRepo := repository.NewMongoChatRepository(collection)
	chatUsecase := usecases.NewChatUsecase(usecases.Deps{
		ChatRepo: chatMongoRepo,
	})

	srv, err := server.NewChatServer(resourcesShutdownCtx, server.Deps{
		ChatUsecase: chatUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

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
	grpcServerOptions := server.UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterChatServiceServer(grpcServer, srv)

	if err = app.Run(ctx, grpcServer); err != nil {
		log.Fatalf("run: %v", err)
	}

	log.Print("servers is stopped")
	resourcesShutdownCtxCancel()
	log.Print("wait shutdown resources")
	time.Sleep(time.Second * 5)

	defer log.Print("app is stopped")
}
