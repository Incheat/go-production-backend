// Package auth defines the controller for the auth API.
package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/incheat/go-playground/services/auth/pkg/model"
)

// Controller is the controller for the auth API.
type Controller struct {
	accessToken      AccessTokenMaker
	refreshToken     RefreshTokenMaker
	refreshTokenRepo RefreshTokenRepository
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

// RefreshTokenRepository is the interface for the refresh token repository.
type RefreshTokenRepository interface {
	SaveRefreshTokenSession(ctx context.Context, session *model.RefreshTokenSession) error
}

// NewController creates a new Controller.
func NewController(accessToken AccessTokenMaker, refreshToken RefreshTokenMaker, refreshTokenRepo RefreshTokenRepository) *Controller {
	return &Controller{accessToken: accessToken, refreshToken: refreshToken, refreshTokenRepo: refreshTokenRepo}
}

// LoginWithEmailAndPassword logs in a user with email and password.
func (c *Controller) LoginWithEmailAndPassword(ctx context.Context, email string, _ string, userAgent, ipAddress string) (*LoginResult, error) {

	memberID := email

	accessToken, err := c.accessToken.CreateToken(memberID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := c.refreshToken.CreateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	maxAge := c.refreshToken.MaxAge()
	refreshEndPoint := c.refreshToken.RefreshEndPoint()

	refreshTokenSession := &model.RefreshTokenSession{
		ID:        uuid.NewString(),
		MemberID:  memberID,
		TokenHash: refreshToken,
		ExpiresAt: now.Add(time.Duration(maxAge) * time.Second),
		CreatedAt: now,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
	err = c.refreshTokenRepo.SaveRefreshTokenSession(ctx, refreshTokenSession)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshMaxAgeSec: maxAge,
		RefreshEndPoint:  refreshEndPoint,
	}, nil
}
