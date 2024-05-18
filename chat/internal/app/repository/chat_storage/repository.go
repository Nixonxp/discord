package repository

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"sync"
)

type InMemoryChatRepository struct {
	storage map[uint64]*models.Message
	mu      sync.RWMutex
}

var _ usecases.ChatsStorage = (*InMemoryChatRepository)(nil)

func NewInMemoryChatRepository() *InMemoryChatRepository {
	return &InMemoryChatRepository{
		storage: make(map[uint64]*models.Message),
	}
}

func (r *InMemoryChatRepository) CreateMessage(_ context.Context, message *models.Message) (*models.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[message.Id] = message
	return r.storage[message.Id], nil
}

func (r *InMemoryChatRepository) GetMessages(_ context.Context, userId uint64, chatId uint64) ([]*models.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	messages := make([]*models.Message, 0)

	for _, v := range r.storage {
		if v.ChatId == chatId && v.UserId == userId {
			messages = append(messages, v)
		}
	}

	return messages, nil
}
