package models

import "time"

type Message struct {
	Id        uint64
	Text      string
	Timestamp time.Time
}

type Messages struct {
	Data []*Message
}

type ActionInfo struct {
	Success bool
}
