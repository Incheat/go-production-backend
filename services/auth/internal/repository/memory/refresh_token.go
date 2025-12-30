// Package memoryrepo defines the memory repository for the refresh token service.
package memoryrepo

import (
	"context"
	"sync"

	"github.com/incheat/go-production-backend/services/auth/internal/repository"
	"github.com/incheat/go-production-backend/services/auth/pkg/model"
)

// RefreshTokenRepository defines a memory refresh token repository.
type RefreshTokenRepository struct {
	sync.RWMutex
	data map[string]*model.RefreshTokenSession
}

// NewRefreshTokenRepository creates a new memory refresh token repository.
func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		data: make(map[string]*model.RefreshTokenSession),
	}
}

// GetRefreshTokenSession gets a refresh token session by token hash.
func (r *RefreshTokenRepository) GetRefreshTokenSession(_ context.Context, refreshToken model.RefreshToken) (*model.RefreshTokenSession, error) {
	r.RLock()
	defer r.RUnlock()
	tokenHash := string(refreshToken)
	refreshTokenSession, ok := r.data[tokenHash]
	if !ok {
		return nil, repository.ErrRefreshTokenNotFound
	}
	return refreshTokenSession, nil
}

// SaveRefreshTokenSession saves a refresh token session.
func (r *RefreshTokenRepository) SaveRefreshTokenSession(_ context.Context, refreshTokenSession *model.RefreshTokenSession) error {
	r.Lock()
	defer r.Unlock()

	tokenHash := string(refreshTokenSession.TokenHash)
	_, ok := r.data[tokenHash]
	if ok {
		return repository.ErrRefreshTokenAlreadyExists
	}

	r.data[tokenHash] = refreshTokenSession
	return nil
}
