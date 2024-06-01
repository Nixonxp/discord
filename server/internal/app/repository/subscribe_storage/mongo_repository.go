package repository

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
	repository "github.com/Nixonxp/discord/server/internal/app/repository/server_storage"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	log "github.com/Nixonxp/discord/server/pkg/logger"
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
		"server_id": subscribe.ServerId,
		"user_id":   subscribe.UserId,
	}

	_, err := r.mongo.UpdateOne(ctx, bson.D{{"_id", uuid.UUID(subscribe.Id)}}, bson.M{
		"$set": addSubscribe,
	}, option)

	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("server_id", subscribe.ServerId).Error("create subscribe error repo")
		return err
	}

	return nil
}

func (r *MongoSubscribeRepository) DeleteSubscribe(ctx context.Context, serverId models.ServerID, userId models.UserID) error {
	filter := bson.M{
		"server_id": serverId,
		"user_id":   userId,
	}

	result, err := r.mongo.DeleteOne(ctx, filter)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("server_id", serverId).Error("delete subscribe error repo")
		return err
	}

	if result.DeletedCount == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *MongoSubscribeRepository) GetByUserId(ctx context.Context, userId models.UserID) ([]*models.SubscribeInfo, error) {
	filter := bson.D{{"user_id", userId}}
	cursor, err := r.mongo.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var subscribes []*models.SubscribeInfo
	err = cursor.All(ctx, &subscribes)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("userId", userId.String()).Error("search subscribe by user error repo")
		return nil, err
	}

	return subscribes, nil
}
