package models

import (
	"github.com/google/uuid"
	"time"
)

type MessageID uuid.UUID

func (v MessageID) String() string {
	return uuid.UUID(v).String()
}

type OwnerID uuid.UUID

func (v OwnerID) String() string {
	return uuid.UUID(v).String()
}

type ChatID uuid.UUID

func (v ChatID) String() string {
	return uuid.UUID(v).String()
}

type UserID uuid.UUID

func (v UserID) String() string {
	return uuid.UUID(v).String()
}

type Chat struct {
	Id       ChatID  `bson:"_id"`
	Type     string  `bson:"type"`
	OwnerId  OwnerID `bson:"owner_id"`
	MetaData string  `bson:"metadata"`
}

type Message struct {
	Id        MessageID `bson:"_id"`
	Text      string    `bson:"text"`
	ChatId    ChatID    `bson:"chat_id"`
	OwnerId   OwnerID   `bson:"owner_id"`
	Timestamp time.Time `bson:"timestamp"`
}

type Messages struct {
	Data []*Message
}

type ActionInfo struct {
	Success bool
}
