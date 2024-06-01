package models

import "github.com/google/uuid"

type ActionInfo struct {
	Success bool
}

type ChannelID uuid.UUID

func (v ChannelID) String() string {
	return uuid.UUID(v).String()
}

type UserID uuid.UUID

func (v UserID) String() string {
	return uuid.UUID(v).String()
}

type Channel struct {
	Id      ChannelID `bson:"_id"`
	Name    string    `bson:"name"`
	OwnerId UserID    `bson:"owner_id"`
}
