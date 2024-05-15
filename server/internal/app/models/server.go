package models

import "time"

type ServerInfo struct {
	Id   uint64
	Name string
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
