package main

import (
	"context"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	"github.com/Nixonxp/discord/chat/internal/app/server"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/chat/internal/middleware/errors"
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

	config := server.Config{
		GRPCPort: ":8080",
		HTTPPort: ":8081",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}
	srv, err := server.NewChatServer(ctx, config, server.Deps{
		ChatUsecase: chatUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	app, err := application.NewApp(srv)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err = app.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
