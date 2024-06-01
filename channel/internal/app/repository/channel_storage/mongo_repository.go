package repository

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/models"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
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
	DeleteOne(ctx context.Context, filter interface{},
		opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type MongoChannelRepository struct {
	mongo MongoCollectionInterface
	log   *log.Logger
}

var _ usecases.ChannelStorage = (*MongoChannelRepository)(nil)

func NewMongoChannelRepository(mongo MongoCollectionInterface, log *log.Logger) *MongoChannelRepository {
	return &MongoChannelRepository{
		mongo: mongo,
		log:   log,
	}
}

func (r *MongoChannelRepository) CreateChannel(ctx context.Context, channel models.Channel) error {
	upsert := true

	option := &options.UpdateOptions{}
	option.Upsert = &upsert

	addMessage := bson.M{
		"name":     channel.Name,
		"owner_id": channel.OwnerId,
	}

	_, err := r.mongo.UpdateOne(ctx, bson.D{{"_id", uuid.UUID(channel.Id)}}, bson.M{
		"$set": addMessage,
	}, option)

	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("name", channel.Name).Error("create channel error repo")
		return err
	}

	return nil
}

const notFoundErrorStr = "mongo: no documents in result"

func (r *MongoChannelRepository) GetChannelById(ctx context.Context, id models.ChannelID) (*models.Channel, error) {
	result := r.mongo.FindOne(context.Background(), bson.M{"_id": uuid.MustParse(id.String())})

	channel := &models.Channel{}
	err := result.Decode(channel)
	if err != nil {
		if err.Error() == notFoundErrorStr {
			return nil, models.ErrNotFound
		}

		r.log.WithContext(ctx).WithError(err).WithField("id", id).Error("decode channel struct error repo")
		return nil, err
	}

	return channel, nil
}

func (r *MongoChannelRepository) DeleteChannel(ctx context.Context, channelId models.ChannelID) error {
	filter := bson.M{
		"_id": channelId,
	}

	result, err := r.mongo.DeleteOne(ctx, filter)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("channelId", channelId.String()).Error("delete channel error repo")
		return err
	}

	if result.DeletedCount == 0 {
		return models.ErrNotFound
	}

	return nil
}
