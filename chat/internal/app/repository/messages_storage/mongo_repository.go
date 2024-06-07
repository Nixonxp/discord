package repository

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoCollectionInterface interface {
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
}

type MongoMessagesRepository struct {
	mongo MongoCollectionInterface
}

var _ usecases.MessagesStorage = (*MongoMessagesRepository)(nil)

func NewMongoMessagesRepository(mongo MongoCollectionInterface) *MongoMessagesRepository {
	return &MongoMessagesRepository{
		mongo: mongo,
	}
}

func (r *MongoMessagesRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	upsert := true

	option := &options.UpdateOptions{}
	option.Upsert = &upsert

	addMessage := bson.M{
		"text":      message.Text,
		"chat_id":   uuid.UUID(message.ChatId),
		"owner_id":  uuid.UUID(message.OwnerId),
		"timestamp": primitive.Timestamp{T: uint32(time.Now().Unix())},
	}

	_, err := r.mongo.UpdateOne(ctx, bson.D{{"_id", uuid.UUID(message.Id)}}, bson.M{
		"$set": addMessage,
	}, option)

	if err != nil {
		return err
	}

	return nil
}

func (r *MongoMessagesRepository) GetMessages(ctx context.Context, chatId models.ChatID) ([]*models.Message, error) {
	filter := bson.D{{"chat_id", uuid.UUID(chatId)}}
	cursor, err := r.mongo.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var messages []*models.Message
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
