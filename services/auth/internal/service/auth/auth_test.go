package authservice_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	usergateway "github.com/incheat/go-production-backend/services/auth/internal/gateway/user"
	authservice "github.com/incheat/go-production-backend/services/auth/internal/service/auth"
	"github.com/incheat/go-production-backend/services/auth/pkg/model"
	usermodel "github.com/incheat/go-production-backend/services/user/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockAccessTokenMaker struct {
	mock.Mock
}

func (m *MockAccessTokenMaker) CreateToken(id string) (model.AccessToken, error) {
	args := m.Called(id)
	return args.Get(0).(model.AccessToken), args.Error(1)
}

type MockRefreshTokenMaker struct {
	mock.Mock
}

func (m *MockRefreshTokenMaker) CreateToken() (model.RefreshToken, error) {
	args := m.Called()
	return args.Get(0).(model.RefreshToken), args.Error(1)
}
func (m *MockRefreshTokenMaker) MaxAge() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockRefreshTokenMaker) RefreshEndPoint() string {
	args := m.Called()
	return args.String(0)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) SaveRefreshTokenSession(ctx context.Context, session *model.RefreshTokenSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

type MockUserGateway struct {
	mock.Mock
}

func (m *MockUserGateway) VerifyCredentials(ctx context.Context, email string, password string) (usermodel.User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(usermodel.User), args.Error(1)
}

func TestUnitLoginWithEmailAndPassword_Success(t *testing.T) {
	ctx := context.Background()
	id := "123"
	email := "user@example.com"
	userAgent := "test-agent"
	ip := "127.0.0.1"
	accessToken := model.AccessToken("access-token")
	refreshToken := model.RefreshToken("refresh-token")
	maxAge := 3600
	endpoint := "/auth/refresh"
	user := usermodel.User{ID: id, Email: email}

	accessMock := new(MockAccessTokenMaker)
	refreshMock := new(MockRefreshTokenMaker)
	repoMock := new(MockRefreshTokenRepository)
	userGatewayMock := new(MockUserGateway)

	userGatewayMock.On("VerifyCredentials", mock.Anything, email, "password-hash").Return(user, nil).Once()
	accessMock.On("CreateToken", email).Return(accessToken, nil).Once()
	refreshMock.On("CreateToken").Return(refreshToken, nil).Once()
	refreshMock.On("MaxAge").Return(maxAge)
	refreshMock.On("RefreshEndPoint").Return(endpoint)

	repoMock.On("SaveRefreshTokenSession", mock.Anything, mock.MatchedBy(func(sess *model.RefreshTokenSession) bool {
		assert.Equal(t, email, sess.MemberID)
		assert.NotEmpty(t, sess.ID)
		return true
	})).Return(nil).Once()

	ctrl := authservice.New(accessMock, refreshMock, repoMock, userGatewayMock)
	result, err := ctrl.LoginWithEmailAndPassword(ctx, email, "password-hash", userAgent, ip)

	require.NoError(t, err)
	assert.Equal(t, accessToken, result.AccessToken)
	assert.Equal(t, refreshToken, result.RefreshToken)
}

func TestUnitLoginWithEmailAndPassword_Errors(t *testing.T) {
	ctx := context.Background()
	id := "123"
	email := "user@example.com"

	tests := []struct {
		name       string
		setupMocks func(a *MockAccessTokenMaker, r *MockRefreshTokenMaker, repo *MockRefreshTokenRepository, ug *MockUserGateway)
		checkErr   func(t *testing.T, err error)
	}{
		{
			name: "credentials error - user not found",
			setupMocks: func(_ *MockAccessTokenMaker, _ *MockRefreshTokenMaker, _ *MockRefreshTokenRepository, ug *MockUserGateway) {
				ug.On("VerifyCredentials", mock.Anything, email, "password-hash").
					Return(usermodel.User{}, usergateway.ErrUserNotFound).Once()
			},
			checkErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, authservice.ErrInvalidCredentials)
				assert.ErrorContains(t, err, "user not found")
			},
		},
		{
			name: "access token generation error",
			setupMocks: func(a *MockAccessTokenMaker, _ *MockRefreshTokenMaker, _ *MockRefreshTokenRepository, ug *MockUserGateway) {
				ug.On("VerifyCredentials", mock.Anything, email, "password-hash").
					Return(usermodel.User{ID: id, Email: email}, nil).Once()
				a.On("CreateToken", email).
					Return(model.AccessToken(""), errors.New("maker-fail")).Once()
			},
			checkErr: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "maker-fail")
			},
		},
		{
			name: "save session error (with dynamic ID description)",
			setupMocks: func(a *MockAccessTokenMaker, r *MockRefreshTokenMaker, repo *MockRefreshTokenRepository, ug *MockUserGateway) {
				ug.On("VerifyCredentials", mock.Anything, email, "password-hash").Return(usermodel.User{ID: id, Email: email}, nil).Once()
				a.On("CreateToken", email).Return(model.AccessToken("at"), nil).Once()
				r.On("CreateToken").Return(model.RefreshToken("rt"), nil).Once()
				r.On("MaxAge").Return(3600)
				r.On("RefreshEndPoint").Return("/ref")

				repo.On("SaveRefreshTokenSession", mock.Anything, mock.Anything).
					Return(errors.New("db-down")).Once()
			},
			checkErr: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "persist refresh token session")
				assert.ErrorContains(t, err, "db-down")
			},
		},
		{
			name: "gateway timeout error",
			setupMocks: func(_ *MockAccessTokenMaker, _ *MockRefreshTokenMaker, _ *MockRefreshTokenRepository, ug *MockUserGateway) {
				timeoutErr := fmt.Errorf("rpc error: %w", context.DeadlineExceeded)
				ug.On("VerifyCredentials", mock.Anything, email, "password-hash").
					Return(usermodel.User{}, timeoutErr).Once()
			},
			checkErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, context.DeadlineExceeded)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, r, repo, ug := new(MockAccessTokenMaker), new(MockRefreshTokenMaker), new(MockRefreshTokenRepository), new(MockUserGateway)
			tt.setupMocks(a, r, repo, ug)

			ctrl := authservice.New(a, r, repo, ug)
			_, err := ctrl.LoginWithEmailAndPassword(ctx, email, "password-hash", "agent", "ip")

			require.Error(t, err)
			tt.checkErr(t, err)

			a.AssertExpectations(t)
			r.AssertExpectations(t)
			repo.AssertExpectations(t)
			ug.AssertExpectations(t)
		})
	}
}
