// Package userservice defines the service for the user service.
package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/incheat/go-production-backend/services/user/internal/repository"
	"github.com/incheat/go-production-backend/services/user/pkg/model"
)

var (
	// ErrInvalidUserCredentials is returned when invalid credentials are provided, e.g. wrong password.
	ErrInvalidUserCredentials = errors.New("invalid user credentials")

	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
)

// Service is the controller for the auth API.
type Service struct {
	userRepo Repository
}

// Repository is the interface for the member repository.
type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, email string, user *model.User) error
}

// New creates a new Service.
func New(userRepo Repository) *Service {
	return &Service{userRepo: userRepo}
}

// VerifyUserCredentials verifies a user's credentials.
func (s *Service) VerifyUserCredentials(ctx context.Context, email string, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		default:
			return nil, fmt.Errorf("get user by email: %w", err)
		}
	}
	if user.PasswordHash != password {
		return nil, ErrInvalidUserCredentials
	}
	return user, nil
}
