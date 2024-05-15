package repository

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/models"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
)

type InMemoryChannelRepository struct {
	storage map[uint64]*models.Channel
}

var _ usecases.ChannelStorage = (*InMemoryChannelRepository)(nil)

func NewInMemoryChannelRepository() *InMemoryChannelRepository {
	return &InMemoryChannelRepository{
		storage: make(map[uint64]*models.Channel),
	}
}

func (r *InMemoryChannelRepository) CreateChannel(_ context.Context, channel models.Channel) (*models.Channel, error) {
	for _, v := range r.storage {
		if v.Name == channel.Name {
			return &models.Channel{}, models.ErrAlreadyExists
		}
	}

	r.storage[channel.Id] = &channel

	return &channel, nil
}
