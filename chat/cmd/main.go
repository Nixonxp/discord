package main

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/queue"
	chat_repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/messages_storage"
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
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"os/signal"
	"syscall"
	"time"
)

const MongoCollectionMessages = "messages"
const MongoCollectionChat = "chat"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	resourcesShutdownCtx, resourcesShutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer resourcesShutdownCtxCancel()

	config := application.Config{
		GRPCPort: ":8580", // todo return
		//HTTPPort:  ":8081",
		//DebugPort: ":8084",
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

	/// KAFKA
	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", "messages", 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	collection, err := mongoCollection.NewCollection(resourcesShutdownCtx,
		MongoCollectionMessages,
		&mongoCollection.Config{
			MongoHost:     "localhost", // todo return
			MongoDb:       "discord",
			MongoPort:     "27117", // todo return
			MongoUser:     "discord",
			MongoPassword: "example", // todo env
		},
	)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	defer collection.DisconnectClient()

	messagesMongoRepo := repository.NewMongoMessagesRepository(collection)

	chatCollection, err := collection.NewCollection(MongoCollectionChat)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}

	chatMongoRepo := chat_repository.NewMongoChatRepository(chatCollection)
	chatUsecase := usecases.NewChatUsecase(usecases.Deps{
		MessagesRepo: messagesMongoRepo,
		ChatRepo:     chatMongoRepo,
		KafkaConn:    conn,
	})

	queueUsecase := usecases.NewQueueUsecase(messagesMongoRepo)
	queueHandler := queue.NewQueue(queue.Deps{QueueUsecase: queueUsecase})
	go func() {
		err := queueHandler.Run(ctx)
		if err != nil {
			return
		}
	}()

	srv, err := server.NewChatServer(resourcesShutdownCtx, server.Deps{
		ChatUsecase: chatUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	globalLimiter := rate_limiter.NewRateLimiter(10000)
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
