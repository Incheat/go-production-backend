// Package userrepo defines the memory repository for the member service.
package userrepo

import (
	"context"
	"sync"

	"github.com/incheat/go-production-backend/services/user/internal/repository"
	"github.com/incheat/go-production-backend/services/user/pkg/model"
)

// UserRepository defines a memory user repository.
type UserRepository struct {
	sync.RWMutex
	data map[string]*model.User
}

// NewUserRepository creates a new memory user repository.
func NewUserRepository() *UserRepository {
	user := &model.User{
		ID:           "1",
		Email:        "test@example.com",
		PasswordHash: "password",
	}
	return &UserRepository{
		data: map[string]*model.User{
			user.Email: user,
		},
	}
}

// GetUserByEmail gets a user by email.
func (r *UserRepository) GetUserByEmail(_ context.Context, email string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()
	user, ok := r.data[email]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

// CreateUser creates a new user.
func (r *UserRepository) CreateUser(_ context.Context, email string, user *model.User) error {
	r.Lock()
	defer r.Unlock()
	_, ok := r.data[email]
	if ok {
		return repository.ErrUserAlreadyExists
	}

	r.data[email] = user
	return nil
}
