package repository

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	"sync"
)

type InMemoryUserRepository struct {
	storage map[string]*models.User
	mu      sync.RWMutex
}

var _ usecases.UsersStorage = (*InMemoryUserRepository)(nil)

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		storage: make(map[string]*models.User),
	}
}

func (r *InMemoryUserRepository) CreateUser(_ context.Context, user *models.User) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.storage {
		if v.Login == user.Login {
			return &models.User{}, models.ErrAlreadyExists
		}
	}

	r.storage[user.UserID.String()] = user

	return user, nil
}

func (r *InMemoryUserRepository) LoginUser(_ context.Context, loginInfo *models.Login) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for k, v := range r.storage {
		if v.Login == loginInfo.Login {
			if v.Password == loginInfo.Password {
				return r.storage[k], nil
			}
			break
		}
	}

	return &models.User{}, models.ErrCredInvalid
}
