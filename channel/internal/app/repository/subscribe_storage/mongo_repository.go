package repository

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/models"
	repository "github.com/Nixonxp/discord/channel/internal/app/repository/channel_storage"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSubscribeRepository struct {
	mongo repository.MongoCollectionInterface
	log   *log.Logger
}

var _ usecases.SubscribeStorage = (*MongoSubscribeRepository)(nil)

func NewMongoSubscribeRepository(mongo repository.MongoCollectionInterface, log *log.Logger) *MongoSubscribeRepository {
	return &MongoSubscribeRepository{
		mongo: mongo,
		log:   log,
	}
}

func (r *MongoSubscribeRepository) CreateSubscribe(ctx context.Context, subscribe models.SubscribeInfo) error {
	upsert := true

	option := &options.UpdateOptions{}
	option.Upsert = &upsert

	addSubscribe := bson.M{
		"channel_id": subscribe.ChannelId,
		"user_id":    subscribe.UserId,
	}

	_, err := r.mongo.UpdateOne(ctx, bson.D{{"_id", uuid.UUID(subscribe.Id)}}, bson.M{
		"$set": addSubscribe,
	}, option)

	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("channel_id", subscribe.ChannelId).Error("create subscribe error repo")
		return err
	}

	return nil
}

func (r *MongoSubscribeRepository) DeleteSubscribe(ctx context.Context, channelId models.ChannelID, userId models.UserID) error {
	filter := bson.M{
		"channel_id": channelId,
		"user_id":    userId,
	}

	result, err := r.mongo.DeleteOne(ctx, filter)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("channel_id", channelId).Error("delete subscribe error repo")
		return err
	}

	if result.DeletedCount == 0 {
		return models.ErrNotFound
	}

	return nil
}
