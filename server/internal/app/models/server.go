package models

import "time"
import "github.com/google/uuid"

type ServerID uuid.UUID

func (v ServerID) String() string {
	return uuid.UUID(v).String()
}

type UserID uuid.UUID

func (v UserID) String() string {
	return uuid.UUID(v).String()
}

type ServerInfo struct {
	Id      ServerID `bson:"_id"`
	Name    string   `bson:"name"`
	OwnerId UserID   `bson:"owner_id"`
}

type ActionInfo struct {
	Success bool
}

type SearchServerByUserIdInfo struct {
	Success bool
}

type GetMessagesInfo struct {
	Messages []*Message
}

type Message struct {
	Id        uint64
	Text      string
	Timestamp time.Time
}
