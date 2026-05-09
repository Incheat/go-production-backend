// Package authservice defines the service for the auth API.
package authservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	usergateway "github.com/incheat/go-production-backend/services/auth/internal/gateway/user"
	"github.com/incheat/go-production-backend/services/auth/pkg/model"
	usermodel "github.com/incheat/go-production-backend/services/user/pkg/model"
)

var (
	// ErrInvalidCredentials is the error for when invalid credentials are provided.
	ErrInvalidCredentials = errors.New("invalid credentials")
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
	VerifyCredentials(ctx context.Context, email string, password string) (usermodel.User, error)
}

// New creates a new Service.
func New(accessToken AccessTokenMaker, refreshToken RefreshTokenMaker, refreshTokenRepo RefreshTokenRepository, userGateway UserGateway) *Service {
	return &Service{accessToken: accessToken, refreshToken: refreshToken, refreshTokenRepo: refreshTokenRepo, userGateway: userGateway}
}

// LoginWithEmailAndPassword logs in a user with email and password.
func (s *Service) LoginWithEmailAndPassword(ctx context.Context, email string, password string, userAgent, ipAddress string) (*LoginResult, error) {

	user, err := s.userGateway.VerifyCredentials(ctx, email, password)
	if err != nil {
		switch {
		// business errors: precisely capture those that should be viewed as 401
		case errors.Is(err, usergateway.ErrUserNotFound),
			errors.Is(err, usergateway.ErrInvalidPassword):
			return nil, fmt.Errorf("%w: %v", ErrInvalidCredentials, err)

		// special cases: client cancel or timeout (optional, but strictly projects usually keep it)
		case errors.Is(err, context.Canceled),
			errors.Is(err, context.DeadlineExceeded):
			return nil, err // original error, let Handler decide how to handle it

		// technical errors: all unexpected Gateway failures
		default:
			return nil, fmt.Errorf("verify credentials via gateway: %w", err)
		}
	}

	memberID := user.Email

	accessToken, err := s.accessToken.CreateToken(memberID)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := s.refreshToken.CreateToken()
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
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
		return nil, fmt.Errorf("persist refresh token session (ID: %s): %w", refreshTokenSession.ID, err)
	}

	return &LoginResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshMaxAgeSec: maxAge,
		RefreshEndPoint:  refreshEndPoint,
	}, nil
}
