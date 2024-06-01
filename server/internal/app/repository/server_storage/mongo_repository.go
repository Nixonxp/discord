package repository

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	log "github.com/Nixonxp/discord/server/pkg/logger"
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

type MongoServerRepository struct {
	mongo MongoCollectionInterface
	log   *log.Logger
}

var _ usecases.ServerStorage = (*MongoServerRepository)(nil)

func NewMongoServerRepository(mongo MongoCollectionInterface, log *log.Logger) *MongoServerRepository {
	return &MongoServerRepository{
		mongo: mongo,
		log:   log,
	}
}

func (r *MongoServerRepository) CreateServer(ctx context.Context, server models.ServerInfo) error {
	upsert := true

	option := &options.UpdateOptions{}
	option.Upsert = &upsert

	addMessage := bson.M{
		"name":     server.Name,
		"owner_id": server.OwnerId,
	}

	_, err := r.mongo.UpdateOne(ctx, bson.D{{"_id", uuid.UUID(server.Id)}}, bson.M{
		"$set": addMessage,
	}, option)

	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("name", server.Name).Error("create server error repo")
		return err
	}

	return nil
}

func (r *MongoServerRepository) SearchServers(ctx context.Context, serverName string) ([]*models.ServerInfo, error) {
	filter := bson.D{{"name", serverName}}
	cursor, err := r.mongo.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var servers []*models.ServerInfo
	err = cursor.All(ctx, &servers)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("serverName", serverName).Error("search server error repo")
		return nil, err
	}

	return servers, nil
}

const notFoundErrorStr = "mongo: no documents in result"

func (r *MongoServerRepository) GetServerById(ctx context.Context, id string) (*models.ServerInfo, error) {
	result := r.mongo.FindOne(context.Background(), bson.M{"_id": uuid.MustParse(id)})

	server := &models.ServerInfo{}
	err := result.Decode(server)
	if err != nil {
		if err.Error() == notFoundErrorStr {
			return nil, models.ErrNotFound
		}

		r.log.WithContext(ctx).WithError(err).WithField("id", id).Error("decode server struct error repo")
		return nil, err
	}

	return server, nil
}
