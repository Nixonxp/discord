package mongo

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Collection struct {
	collection  *mongo.Collection
	database    *mongo.Database
	clientMongo *mongo.Client
}

type Config struct {
	MongoDb       string
	MongoHost     string
	MongoPort     string
	MongoUser     string
	MongoPassword string
}

func NewCollection(ctx context.Context, collectionName string, config *Config) (*Collection, error) {
	mongoUrl := "mongodb://" + config.MongoHost + ":" + config.MongoPort
	credentials := options.Credential{
		Username:   config.MongoUser,
		Password:   config.MongoPassword,
		AuthSource: config.MongoDb,
	}

	clientMongo, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl).SetAuth(credentials))
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := clientMongo.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	database := clientMongo.Database(config.MongoDb)

	return &Collection{
		collection:  database.Collection(collectionName),
		database:    database,
		clientMongo: clientMongo,
	}, nil
}

func (c *Collection) DisconnectClient() error {
	if err := c.clientMongo.Disconnect(context.Background()); err != nil {
		return err
	}

	return nil
}

func (c *Collection) Clone() (*Collection, error) {
	newCollection, err := c.collection.Clone()
	if err != nil {
		return nil, err
	}

	return &Collection{
		collection:  newCollection,
		clientMongo: c.clientMongo,
	}, nil
}

func (c *Collection) NewCollection(collectionName string) (*Collection, error) {
	return &Collection{
		collection:  c.database.Collection(collectionName),
		database:    c.database,
		clientMongo: c.clientMongo,
	}, nil
}

func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongodb.UpdateOne")
	defer span.Finish()

	span.LogFields(
		log.Object("update", update),
	)
	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongodb.Find")
	defer span.Finish()

	return c.collection.Find(ctx, filter, opts...)
}

func (c *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongodb.FindOne")
	defer span.Finish()

	return c.collection.FindOne(ctx, filter, opts...)
}

func (c *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongodb.DeleteOne")
	defer span.Finish()

	return c.collection.DeleteOne(ctx, filter, opts...)
}
