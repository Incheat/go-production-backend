// Package authservice defines the service for the auth API.
package authservice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/incheat/go-playground/services/auth/pkg/model"
	usermodel "github.com/incheat/go-playground/services/user/pkg/model"
)

// Service is the service for the auth API.
type Service struct {
	accessToken      AccessTokenMaker
	refreshToken     RefreshTokenMaker
	refreshTokenRepo RefreshTokenRepository
	userGateway      UserGateway
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

// UserGateway is the interface for the user gateway.
type UserGateway interface {
	VerifyCredentials(ctx context.Context, email string, password string) (*usermodel.User, error)
}

// New creates a new Service.
func New(accessToken AccessTokenMaker, refreshToken RefreshTokenMaker, refreshTokenRepo RefreshTokenRepository, userGateway UserGateway) *Service {
	return &Service{accessToken: accessToken, refreshToken: refreshToken, refreshTokenRepo: refreshTokenRepo, userGateway: userGateway}
}

// LoginWithEmailAndPassword logs in a user with email and password.
func (s *Service) LoginWithEmailAndPassword(ctx context.Context, email string, password string, userAgent, ipAddress string) (*LoginResult, error) {

	user, err := s.userGateway.VerifyCredentials(ctx, email, password)
	if err != nil {
		return nil, err
	}

	memberID := user.Email

	accessToken, err := s.accessToken.CreateToken(memberID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshToken.CreateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	maxAge := s.refreshToken.MaxAge()
	refreshEndPoint := s.refreshToken.RefreshEndPoint()

	refreshTokenSession := &model.RefreshTokenSession{
		ID:        uuid.NewString(),
		MemberID:  memberID,
		TokenHash: refreshToken,
		ExpiresAt: now.Add(time.Duration(maxAge) * time.Second),
		CreatedAt: now,
		RevokedAt: time.Time{}, // not revoked yet, set to zero value
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
	err = s.refreshTokenRepo.SaveRefreshTokenSession(ctx, refreshTokenSession)
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
