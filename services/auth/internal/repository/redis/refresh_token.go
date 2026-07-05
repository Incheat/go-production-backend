// Package redisrepo defines the Redis refresh token repository.
package redisrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/incheat/go-production-backend/services/auth/internal/constant"
	"github.com/incheat/go-production-backend/services/auth/internal/repository"
	"github.com/incheat/go-production-backend/services/auth/pkg/model"
	"github.com/redis/go-redis/v9"
)

// RefreshTokenRepository defines a Redis refresh token repository.
type RefreshTokenRepository struct {
	rdb    *redis.Client
	prefix string
}

// NewRefreshTokenRepository creates a new Redis refresh token repository.
func NewRefreshTokenRepository(rdb *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		rdb:    rdb,
		prefix: constant.RedisRefreshTokenPrefix, // key prefix in Redis
	}
}

// key builds a Redis key from a token hash.
func (r *RefreshTokenRepository) key(hash string) string {
	return r.prefix + hash
}

// GetRefreshTokenSession gets a refresh token session by token hash.
func (r *RefreshTokenRepository) GetRefreshTokenSession(ctx context.Context, refreshToken model.RefreshToken) (*model.RefreshTokenSession, error) {
	tokenHash := string(refreshToken)
	key := r.key(tokenHash)

	data, err := r.rdb.Get(ctx, key).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	var session model.RefreshTokenSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// SaveRefreshTokenSession saves a refresh token session.
func (r *RefreshTokenRepository) SaveRefreshTokenSession(ctx context.Context, session *model.RefreshTokenSession) error {
	tokenHash := string(session.TokenHash)
	key := r.key(tokenHash)

	// Check existence first
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("check if refresh token session exists: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("%w (key: %s)", repository.ErrAlreadyExists, key)
	}

	// Marshal to JSON
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal refresh token session to JSON: %w", err)
	}

	// TTL = ExpiresAt - Now
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Minute // fallback TTL just in case
	}

	// Store session with TTL
	if err := r.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
