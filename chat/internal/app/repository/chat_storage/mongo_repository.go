package repository

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionInterface interface {
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{},
		opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	FindOne(ctx context.Context, filter interface{},
		opts ...*options.FindOneOptions) *mongo.SingleResult
}

type MongoChatRepository struct {
	mongo MongoCollectionInterface
}

const notFoundErrorStr = "mongo: no documents in result"

var _ usecases.ChatStorage = (*MongoChatRepository)(nil)

func NewMongoChatRepository(mongo MongoCollectionInterface) *MongoChatRepository {
	return &MongoChatRepository{
		mongo: mongo,
	}
}

func (r *MongoChatRepository) CreateChat(ctx context.Context, chat *models.Chat) error {
	upsert := true

	option := &options.UpdateOptions{}
	option.Upsert = &upsert

	addMessage := bson.M{
		"type":     chat.Type,
		"owner_id": uuid.UUID(chat.OwnerId),
		"metadata": chat.MetaData,
	}

	_, err := r.mongo.UpdateOne(ctx, bson.D{{"_id", uuid.UUID(chat.Id)}}, bson.M{
		"$set": addMessage,
	}, option)

	if err != nil {
		return err
	}

	return nil
}

func (r *MongoChatRepository) GetChatById(ctx context.Context, chatId models.ChatID) (*models.Chat, error) {
	result := r.mongo.FindOne(context.Background(), bson.M{"_id": uuid.MustParse(chatId.String())})

	chat := &models.Chat{}
	err := result.Decode(chat)
	if err != nil {
		if err.Error() == notFoundErrorStr {
			return nil, models.ErrNotFound
		}

		//r.log.WithContext(ctx).WithError(err).WithField("id", id).Error("decode channel struct error repo")
		return nil, err
	}

	return chat, nil
}

/*func (r *MongoChatRepository) GetChatByOwnerIdAndType(ctx context.Context, ownerId models.OwnerID, chatType string) (*models.Chat, error) { // todo delete?
	result := r.mongo.FindOne(context.Background(), bson.M{"owner_id": uuid.MustParse(ownerId.String()), "type": uuid.MustParse(chatType)})

	chat := &models.Chat{}
	err := result.Decode(chat)
	if err != nil {
		if err.Error() == notFoundErrorStr {
			return nil, models.ErrNotFound
		}

		//r.log.WithContext(ctx).WithError(err).WithField("id", id).Error("decode channel struct error repo")
		return nil, err
	}

	return chat, nil
}*/

func (r *MongoChatRepository) GetChatByMetadataAndType(ctx context.Context, metadata string, chatType string) (*models.Chat, error) {
	result := r.mongo.FindOne(context.Background(), bson.M{"metadata": metadata, "type": chatType})

	chat := &models.Chat{}
	err := result.Decode(chat)
	if err != nil {
		if err.Error() == notFoundErrorStr {
			return nil, models.ErrNotFound
		}

		//r.log.WithContext(ctx).WithError(err).WithField("id", id).Error("decode channel struct error repo")
		return nil, err
	}

	return chat, nil
}
