package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/chat/internal/app/queue"
	chat_repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/messages_storage"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/chat/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/chat/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/chat/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/chat/pkg/api/v1"
	"github.com/Nixonxp/discord/chat/pkg/rate_limiter"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors []grpc.UnaryServerInterceptor
	UnaryInterceptors      []grpc.UnaryServerInterceptor
}

// Deps - server deps
type Deps struct {
	ChatUsecase usecases.UsecaseInterface
}

const MongoCollectionChat = "chat"

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	Deps

	validator *protovalidate.Validator

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	http struct {
		lis  *echo.Echo
		port string
	}

	grpcServer *grpc.Server
}

func NewChatServer(ctx context.Context, s *MainServer) (*ChatServer, error) {
	srv := &ChatServer{}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.SendUserPrivateMessageRequest{},
				&pb.GetUserPrivateMessagesRequest{},
				&pb.CreatePrivateChatRequest{},
				&pb.SendServerMessageRequest{},
				&pb.GetServerMessagesRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}

	mainCollection := s.mongo.GetInstance()
	messagesMongoRepo := repository.NewMongoMessagesRepository(mainCollection)

	chatCollection, err := mainCollection.NewCollection(MongoCollectionChat)
	if err != nil {
		return nil, fmt.Errorf("failed to connect mongo: %v", err)
	}

	chatMongoRepo := chat_repository.NewMongoChatRepository(chatCollection)
	chatUsecase := usecases.NewChatUsecase(usecases.Deps{
		MessagesRepo: messagesMongoRepo,
		ChatRepo:     chatMongoRepo,
		KafkaConn:    s.kafka.GetInstance(),
	})

	queueUsecase := usecases.NewQueueUsecase(messagesMongoRepo)
	queueHandler := queue.NewQueue(queue.Deps{QueueUsecase: queueUsecase})
	go func() {
		err := queueHandler.Run(ctx)
		if err != nil {
			return
		}
	}()

	srv.ChatUsecase = chatUsecase

	globalLimiter := rate_limiter.NewRateLimiter(10000)
	grpcConfig := Config{
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
	grpcServerOptions := UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterChatServiceServer(grpcServer, srv)

	srv.grpcServer = grpcServer

	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
