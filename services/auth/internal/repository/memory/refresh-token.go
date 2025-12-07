// Package memory defines the memory repository for the refresh token service.
package memory

import (
	"context"
	"sync"

	"github.com/incheat/go-playground/services/auth/internal/repository"
	"github.com/incheat/go-playground/services/auth/pkg/model"
)

// RefreshTokenRepository defines a memory refresh token repository.
type RefreshTokenRepository struct {
	sync.RWMutex
	data map[string]*model.RefreshToken
}

// NewRefreshTokenRepository creates a new memory refresh token repository.
func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		data: make(map[string]*model.RefreshToken),
	}
}

// GetRefreshTokenByID gets a refresh token by ID.
func (r *RefreshTokenRepository) GetRefreshTokenByID(_ context.Context, id string) (*model.RefreshToken, error) {
	r.RLock()
	defer r.RUnlock()
	refreshToken, ok := r.data[id]
	if !ok {
		return nil, repository.ErrRefreshTokenNotFound
	}
	return refreshToken, nil
}

// CreateRefreshToken creates a new refresh token.
func (r *RefreshTokenRepository) CreateRefreshToken(_ context.Context, id string, refreshToken *model.RefreshToken) error {
	r.RLock()
	defer r.RUnlock()
	_, ok := r.data[id]
	if ok {
		return repository.ErrRefreshTokenAlreadyExists
	}

	r.Lock()
	defer r.Unlock()
	r.data[id] = refreshToken
	return nil
}
