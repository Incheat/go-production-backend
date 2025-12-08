package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/incheat/go-playground/services/auth/internal/controller/auth"
	"github.com/incheat/go-playground/services/auth/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Testify mocks ---

type MockAccessTokenMaker struct {
	mock.Mock
}

func (m *MockAccessTokenMaker) CreateToken(id string) (model.AccessToken, error) {
	args := m.Called(id)
	return args.Get(0).(model.AccessToken), args.Error(1)
}

func (m *MockAccessTokenMaker) ParseToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
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

func (m *MockRefreshTokenRepository) SaveRefreshTokenSession(
	ctx context.Context,
	session *model.RefreshTokenSession,
) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func TestLoginWithEmailAndPassword_Success(t *testing.T) {
	ctx := context.Background()
	email := "user@example.com"
	userAgent := "test-agent"
	ip := "127.0.0.1"
	accessToken := model.AccessToken("access-token")
	refreshToken := model.RefreshToken("refresh-token")
	maxAge := 3600
	endpoint := "/auth/refresh"

	accessMock := new(MockAccessTokenMaker)
	refreshMock := new(MockRefreshTokenMaker)
	repoMock := new(MockRefreshTokenRepository)

	// Expectations
	accessMock.
		On("CreateToken", email).
		Return(accessToken, nil).
		Once()

	refreshMock.
		On("CreateToken").
		Return(refreshToken, nil).
		Once()

	refreshMock.
		On("MaxAge").
		Return(maxAge)

	refreshMock.
		On("RefreshEndPoint").
		Return(endpoint)

	// We want to inspect the session passed to the repo:
	repoMock.
		On(
			"SaveRefreshTokenSession",
			mock.Anything,
			mock.MatchedBy(func(sess *model.RefreshTokenSession) bool {
				// Basic field checks
				assert.Equal(t, email, sess.MemberID)
				assert.Equal(t, refreshToken, sess.TokenHash)
				assert.Equal(t, userAgent, sess.UserAgent)
				assert.Equal(t, ip, sess.IPAddress)

				// Check time relationship
				assert.True(t, sess.ExpiresAt.After(sess.CreatedAt))

				expectedDuration := time.Duration(maxAge) * time.Second
				actualDuration := sess.ExpiresAt.Sub(sess.CreatedAt)

				// allow small tolerance
				const tolerance = time.Second
				diff := actualDuration - expectedDuration
				if diff < 0 {
					diff = -diff
				}
				assert.LessOrEqual(t, diff, tolerance)

				// ID should not be empty (uuid string)
				assert.NotEmpty(t, sess.ID)

				return true
			}),
		).
		Return(nil).
		Once()

	ctrl := auth.NewController(accessMock, refreshMock, repoMock)

	result, err := ctrl.LoginWithEmailAndPassword(ctx, email, "password", userAgent, ip)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Check returned result
	assert.Equal(t, accessToken, result.AccessToken)
	assert.Equal(t, refreshToken, result.RefreshToken)
	assert.Equal(t, maxAge, result.RefreshMaxAgeSec)
	assert.Equal(t, endpoint, result.RefreshEndPoint)

	accessMock.AssertExpectations(t)
	refreshMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
}

func TestLoginWithEmailAndPassword_Errors(t *testing.T) {
	ctx := context.Background()
	email := "user@example.com"

	tests := []struct {
		name             string
		setupMocks       func(a *MockAccessTokenMaker, r *MockRefreshTokenMaker, repo *MockRefreshTokenRepository)
		expectedErr      error
		expectRepoCalled bool
	}{
		{
			name: "access token error",
			setupMocks: func(a *MockAccessTokenMaker, _ *MockRefreshTokenMaker, _ *MockRefreshTokenRepository) {
				err := errors.New("access error")
				a.On("CreateToken", email).
					Return(model.AccessToken(""), err).
					Once()
			},
			expectedErr:      errors.New("access error"),
			expectRepoCalled: false,
		},
		{
			name: "refresh token error",
			setupMocks: func(a *MockAccessTokenMaker, r *MockRefreshTokenMaker, _ *MockRefreshTokenRepository) {
				a.On("CreateToken", email).
					Return(model.AccessToken("access-token"), nil).
					Once()

				err := errors.New("refresh error")
				r.On("CreateToken").
					Return(model.RefreshToken(""), err).
					Once()
			},
			expectedErr:      errors.New("refresh error"),
			expectRepoCalled: false,
		},
		{
			name: "save session error",
			setupMocks: func(a *MockAccessTokenMaker, r *MockRefreshTokenMaker, repo *MockRefreshTokenRepository) {
				a.On("CreateToken", email).
					Return(model.AccessToken("access-token"), nil).
					Once()

				r.On("CreateToken").
					Return(model.RefreshToken("refresh-token"), nil).
					Once()

				r.On("MaxAge").
					Return(3600)

				r.On("RefreshEndPoint").
					Return("/refresh")

				err := errors.New("save error")
				repo.On("SaveRefreshTokenSession", mock.Anything, mock.AnythingOfType("*model.RefreshTokenSession")).
					Return(err).
					Once()
			},
			expectedErr:      errors.New("save error"),
			expectRepoCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessMock := new(MockAccessTokenMaker)
			refreshMock := new(MockRefreshTokenMaker)
			repoMock := new(MockRefreshTokenRepository)

			tt.setupMocks(accessMock, refreshMock, repoMock)

			ctrl := auth.NewController(accessMock, refreshMock, repoMock)

			result, err := ctrl.LoginWithEmailAndPassword(ctx, email, "password", "agent", "ip")
			require.Error(t, err)
			assert.Nil(t, result)
			assert.EqualError(t, err, tt.expectedErr.Error())

			accessMock.AssertExpectations(t)
			refreshMock.AssertExpectations(t)
			repoMock.AssertExpectations(t)
		})
	}
}
