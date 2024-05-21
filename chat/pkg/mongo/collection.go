package mongo

import (
	"context"
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

	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

func (c *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	return c.collection.Find(ctx, filter, opts...)
}
