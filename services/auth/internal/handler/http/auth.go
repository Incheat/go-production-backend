package http

import (
	"context"
	"time"

	"github.com/google/uuid"
	gen "github.com/incheat/go-playground/services/auth/internal/api/gen/server"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// _ is a placeholder to ensure that Server implements the StrictServerInterface interface.
var _ gen.StrictServerInterface = (*Handler)(nil)

// Handler is the handler for the Auth API.
type Handler struct {}

// NewServer creates a new Server.
func NewHandler() *Handler {
	return &Handler{}
}

// GetMe is the handler for the GetMe endpoint.
func (h *Handler) GetMe(ctx context.Context, request gen.GetMeRequestObject) (gen.GetMeResponseObject, error) {

	id := openapi_types.UUID(uuid.New())
	email := openapi_types.Email("test@example.com")
	name := "Test User"
	createdAt := time.Now()

	dummyUser := gen.User{
		Id:        &id,
		Email:     &email,
		Name:      &name,
		CreatedAt: &createdAt,
	}

	return gen.GetMe200JSONResponse{
		Body: dummyUser,
		Headers: gen.GetMe200ResponseHeaders{
			VersionId: "v1",
		},
	}, nil
}

// Login is the handler for the Login endpoint.
func (h *Handler) Login(ctx context.Context, request gen.LoginRequestObject) (gen.LoginResponseObject, error) {
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	return gen.Login200JSONResponse{
		Body: gen.AuthResponse{
			AccessToken:  &accessToken,
			RefreshToken: &refreshToken,
		},
		Headers: gen.Login200ResponseHeaders{
			VersionId: "v1",
		},
	}, nil
}

// Logout is the handler for the Logout endpoint.
func (h *Handler) Logout(ctx context.Context, request gen.LogoutRequestObject) (gen.LogoutResponseObject, error) {
	return gen.Logout204Response{
		Headers: gen.Logout204ResponseHeaders{
			VersionId: "v1",
		},
	}, nil
}

// Register is the handler for the Register endpoint.
func (h *Handler) Register(ctx context.Context, request gen.RegisterRequestObject) (gen.RegisterResponseObject, error) {
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	return gen.Register201JSONResponse{
		Body: gen.AuthResponse{
			AccessToken:  &accessToken,
			RefreshToken: &refreshToken,
		},
		Headers: gen.Register201ResponseHeaders{
			VersionId: "v1",
		},
	}, nil
}