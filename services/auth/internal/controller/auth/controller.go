// Package auth defines the controller for the auth API.
package auth

import (
	"context"
	"time"

	"github.com/incheat/go-playground/services/auth/pkg/model"
)

// Controller is the controller for the auth API.
type Controller struct {
	accessToken  AccessTokenMaker
	refreshToken RefreshTokenMaker
	redis        RedisClient
}

// RedisClient is the interface for the Redis client.
type RedisClient interface {
	Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

// AccessTokenMaker is the interface for the access token maker.
type AccessTokenMaker interface {
	CreateToken(ID string) (model.AccessToken, error)
	ParseToken(token string) (string, error)
}

// RefreshTokenMaker is the interface for the refresh token maker.
type RefreshTokenMaker interface {
	CreateToken() (model.RefreshToken, error)
	MaxAge() int
	RefreshEndPoint() string
}

// NewController creates a new Controller.
func NewController(accessToken AccessTokenMaker, refreshToken RefreshTokenMaker, redis RedisClient) *Controller {
	return &Controller{accessToken: accessToken, refreshToken: refreshToken, redis: redis}
}

// LoginWithEmailAndPassword logs in a user with email and password.
func (c *Controller) LoginWithEmailAndPassword(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// GenerateAccessTokenByID generates a new access token by ID.
func (c *Controller) GenerateAccessTokenByID(ID string) (model.AccessToken, error) {
	accessToken, err := c.accessToken.CreateToken(ID)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// GenerateRefreshToken generates a new refresh token.
func (c *Controller) GenerateRefreshToken() (model.RefreshToken, error) {
	refreshToken, err := c.refreshToken.CreateToken()
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

// RefreshTokenMaxAge returns the maximum age of the refresh token.
func (c *Controller) RefreshTokenMaxAge() int {
	return c.refreshToken.MaxAge()
}

// RefreshTokenRefreshEndPoint returns the refresh end point of the refresh token.
func (c *Controller) RefreshTokenRefreshEndPoint() string {
	return c.refreshToken.RefreshEndPoint()
}
