package repository

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
)

type InMemoryUserRepository struct {
	storage map[uint64]*models.User
}

var _ usecases.UsersStorage = (*InMemoryUserRepository)(nil)

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		storage: make(map[uint64]*models.User),
	}
}

func (r *InMemoryUserRepository) UpdateUser(_ context.Context, user *models.User) (*models.User, error) {
	for _, v := range r.storage {
		if v.Login == user.Login {
			r.storage[user.UserID].UserID = user.UserID
			r.storage[user.UserID].Name = user.Name
			r.storage[user.UserID].Login = user.Login
			r.storage[user.UserID].Email = user.Email
			return r.storage[user.UserID], nil
		}
	}

	return &models.User{}, models.ErrNotFound
}
