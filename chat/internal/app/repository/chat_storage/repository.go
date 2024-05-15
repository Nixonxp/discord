package repository

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
)

type InMemoryChatRepository struct {
	storage map[uint64]*models.Message
}

var _ usecases.ChatsStorage = (*InMemoryChatRepository)(nil)

func NewInMemoryChatRepository() *InMemoryChatRepository {
	return &InMemoryChatRepository{
		storage: make(map[uint64]*models.Message),
	}
}

func (r *InMemoryChatRepository) CreateMessage(_ context.Context, message models.Message) (*models.Message, error) {
	r.storage[message.Id] = &message
	return r.storage[message.Id], nil
}
