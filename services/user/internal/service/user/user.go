// Package userservice defines the service for the user service.
package userservice

import (
	"context"
	"errors"

	"github.com/incheat/go-playground/services/user/pkg/model"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrUserAlreadyExists is returned when a user already exists.
var ErrUserAlreadyExists = errors.New("user already exists")

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
		return nil, err
	}
	if user.PasswordHash != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
