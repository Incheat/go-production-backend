// Package userrepo defines the memory repository for the member service.
package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-sql-driver/mysql"
	db "github.com/incheat/go-playground/services/user/internal/db/mysql/gen"
	"github.com/incheat/go-playground/services/user/internal/repository"
	"github.com/incheat/go-playground/services/user/pkg/model"
)

// UserRepository defines a memory user repository.
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new user repository.
func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		queries: db.New(dbConn),
	}
}

// GetUserByEmail gets a user by email.
func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	u, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &model.User{
		ID:           strconv.FormatInt(u.ID, 10),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}, nil
}

// CreateUser creates a new user.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	email string,
	user *model.User,
) error {

	_, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			return repository.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
