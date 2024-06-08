package services

import (
	"context"
	"fmt"
	config "github.com/Nixonxp/discord/channel/configs"
	mongoCollection "github.com/Nixonxp/discord/channel/pkg/mongo"
)

type Mongo struct {
	coll *mongoCollection.Collection
}

func (m *Mongo) Init(ctx context.Context, cfg *config.Config) error {
	var err error
	m.coll, err = mongoCollection.NewCollection(ctx,
		cfg.Application.ServiceCollection,
		&mongoCollection.Config{
			MongoHost:     cfg.Application.MongoHost,
			MongoDb:       cfg.Application.MongoDb,
			MongoPort:     cfg.Application.MongoPort,
			MongoUser:     cfg.Application.MongoUser,
			MongoPassword: cfg.Application.MongoPassword,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to connect mongo: %v", err)
	}

	return nil
}

func (m *Mongo) GetInstance() *mongoCollection.Collection {
	return m.coll
}

func (m *Mongo) Ident() string {
	return "mongo"
}

func (m *Mongo) Close(ctx context.Context) error {
	err := m.coll.DisconnectClient()
	if err != nil {
		return err
	}
	return nil
}
