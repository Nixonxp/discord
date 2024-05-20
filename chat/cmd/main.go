package main

import (
	"context"
	repository "github.com/Nixonxp/discord/chat/internal/app/repository/chat_storage"
	"github.com/Nixonxp/discord/chat/internal/app/server"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	middleware "github.com/Nixonxp/discord/chat/internal/middleware/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"log"
	"time"
)

const MongoCollectionMessages = "messages"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := server.Config{
		GRPCPort:      ":8080",
		HTTPPort:      ":8081",
		MongoHost:     "mongodb-chat",
		MongoDb:       "Chat",
		MongoPort:     "27017",
		MongoUser:     "root",
		MongoPassword: "example",
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(),
		},
	}

	// Mongo DB
	ctxMongo, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoUrl := "mongodb://" + config.MongoHost + ":" + config.MongoPort
	credentials := options.Credential{
		Username:   config.MongoUser,
		Password:   config.MongoPassword,
		AuthSource: config.MongoDb,
	}

	clientMongo, err := mongo.Connect(ctxMongo, options.Client().ApplyURI(mongoUrl).SetAuth(credentials))
	if err != nil {
		log.Fatalf("mongo error: %v", err)
	}

	defer func() {
		if err = clientMongo.Disconnect(ctxMongo); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := clientMongo.Ping(ctxMongo, readpref.Primary()); err != nil {
		log.Fatalf("mongo error: %v", err)
	}

	database := clientMongo.Database(config.MongoDb)
	chatMongoRepo := repository.NewMongoChatRepository(database.Collection(MongoCollectionMessages))
	chatUsecase := usecases.NewChatUsecase(usecases.Deps{
		ChatRepo: chatMongoRepo,
	})

	// delivery

	srv, err := server.NewChatServer(ctx, config, server.Deps{
		ChatUsecase: chatUsecase,
	})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err = srv.Run(ctx); err != nil {
		log.Fatalf("run: %v", err)
	}
}
