package models

import "github.com/google/uuid"

type SubscribeID uuid.UUID

func (v SubscribeID) String() string {
	return uuid.UUID(v).String()
}

type SubscribeInfo struct {
	Id       SubscribeID `bson:"_id"`
	ServerId ServerID    `bson:"server_id"`
	UserId   UserID      `bson:"user_id"`
}
