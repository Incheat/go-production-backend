package userservice_test

import (
	"context"
	"errors"
	"testing"

	userservice "github.com/incheat/go-production-backend/services/user/internal/service/user"
	"github.com/incheat/go-production-backend/services/user/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Testify mocks ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, email string, user *model.User) error {
	args := m.Called(ctx, email, user)
	return args.Error(0)
}

// TestUnitVerifyUserCredentials_Success tests the happy path for VerifyUserCredentials.
func TestUnitVerifyUserCredentials_Success(t *testing.T) {
	ctx := context.Background()
	email := "user@example.com"
	password := "password"

	repoMock := new(MockUserRepository)

	expectedUser := &model.User{
		Email:        email,
		PasswordHash: password,
	}

	repoMock.
		On("GetUserByEmail", mock.Anything, email).
		Return(expectedUser, nil).
		Once()

	svc := userservice.New(repoMock)

	got, err := svc.VerifyUserCredentials(ctx, email, password)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, expectedUser, got)

	repoMock.AssertExpectations(t)
}

// TestUnitVerifyUserCredentials_Errors tests the error cases for VerifyUserCredentials.
func TestUnitVerifyUserCredentials_Errors(t *testing.T) {
	ctx := context.Background()
	email := "user@example.com"
	password := "password"

	tests := []struct {
		name       string
		setupMocks func(repo *MockUserRepository)
		wantErr    string
	}{
		{
			name: "repo get user error",
			setupMocks: func(repo *MockUserRepository) {
				repo.
					On("GetUserByEmail", mock.Anything, email).
					Return((*model.User)(nil), errors.New("db error")).
					Once()
			},
			wantErr: "db error",
		},
		{
			name: "invalid credentials",
			setupMocks: func(repo *MockUserRepository) {
				repo.
					On("GetUserByEmail", mock.Anything, email).
					Return(&model.User{
						Email:        email,
						PasswordHash: "different-password",
					}, nil).
					Once()
			},
			wantErr: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(MockUserRepository)
			tt.setupMocks(repoMock)

			svc := userservice.New(repoMock)

			got, err := svc.VerifyUserCredentials(ctx, email, password)
			require.Error(t, err)
			assert.Nil(t, got)
			assert.EqualError(t, err, tt.wantErr)

			repoMock.AssertExpectations(t)
		})
	}
}
