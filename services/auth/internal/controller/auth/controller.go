package auth

import (
	"context"
	"errors"
	"time"

	"github.com/incheat/go-playground/services/auth/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrMemberNotFound = errors.New("member not found")

// ErrMemberAlreadyExists is returned when a member already exists.
var ErrMemberAlreadyExists = errors.New("member already exists")

// Controller is the controller for the auth API.
type Controller struct {
	refreshTokenRepo RefreshTokenRepository
	jwt JWTMaker
	redis RedisClient
}

// RedisClient is the interface for the Redis client.
type RedisClient interface {
	Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

// JWTMaker is the interface for the JWT maker.
type JWTMaker interface {
	CreateToken(userID string) (string, error)
	ParseUserID(tokenStr string) (string, error)
}

// RefreshTokenRepository is the interface for the refresh token repository.
type RefreshTokenRepository interface {
	GetRefreshTokenByID(ctx context.Context, id string) (*model.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, id string, refreshToken *model.RefreshToken) error
}

// NewController creates a new Controller.
func NewController(refreshTokenRepo RefreshTokenRepository, jwt JWTMaker, redis RedisClient) *Controller {
	return &Controller{refreshTokenRepo: refreshTokenRepo, jwt: jwt, redis: redis}
}