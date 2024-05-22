package main

import (
	"context"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	"github.com/Nixonxp/discord/chat/internal/app/server"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/chat/internal/middleware/errors"
	pb "github.com/Nixonxp/discord/chat/pkg/api/v1"
	"github.com/Nixonxp/discord/chat/pkg/application"
	mongoCollection "github.com/Nixonxp/discord/chat/pkg/mongo"
	"google.golang.org/grpc"
	"log"
	"time"
)

const MongoCollectionMessages = "messages"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := application.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
	}

	app, err := application.NewApp(&config)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	ctxMongo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := mongoCollection.NewCollection(ctxMongo,
		MongoCollectionMessages,
		&mongoCollection.Config{
			MongoHost:     "mongodb-chat",
			MongoDb:       "Chat",
			MongoPort:     "27017",
			MongoUser:     "root",
			MongoPassword: "example",
		},
	)
	defer collection.DisconnectClient()

	chatMongoRepo := repository.NewMongoChatRepository(collection)
	chatUsecase := usecases.NewChatUsecase(usecases.Deps{
		ChatRepo: chatMongoRepo,
	})

	srv, err := server.NewChatServer(ctx, server.Deps{
		ChatUsecase: chatUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	grpcConfig := server.Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
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
}
