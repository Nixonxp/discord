package repository

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
)

type InMemoryServerRepository struct {
	storage map[uint64]*models.ServerInfo
}

var _ usecases.ServerStorage = (*InMemoryServerRepository)(nil)

func NewInMemoryServerRepository() *InMemoryServerRepository {
	return &InMemoryServerRepository{
		storage: make(map[uint64]*models.ServerInfo),
	}
}

func (r *InMemoryServerRepository) CreateServer(_ context.Context, server models.ServerInfo) (*models.ServerInfo, error) {
	for _, v := range r.storage {
		if v.Name == server.Name {
			return &models.ServerInfo{}, models.ErrAlreadyExists
		}
	}

	r.storage[server.Id] = &server

	return &server, nil
}
