// Package handler defines the handlers for the Auth API.
package handler

import (
	"context"
	"fmt"

	gen "github.com/incheat/go-playground/services/auth/internal/api/gen/oapi/public/server"
	"github.com/incheat/go-playground/services/auth/internal/constant"
	"github.com/incheat/go-playground/services/auth/internal/controller/auth"
)

// _ is a placeholder to ensure that Server implements the StrictServerInterface interface.
var _ gen.StrictServerInterface = (*Handler)(nil)

// Handler is the handler for the Auth API.
type Handler struct {
	ctrl *auth.Controller
}

// NewHandler creates a new Handler.
func NewHandler(ctrl *auth.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Login is the handler for the Login endpoint.
func (h *Handler) Login(ctx context.Context, request gen.LoginRequestObject) (gen.LoginResponseObject, error) {
	email := string(request.Body.Email)
	password := request.Body.Password
	success, err := h.ctrl.LoginWithEmailAndPassword(ctx, email, password)
	if err != nil {
		return gen.Login500JSONResponse{
			Error: err.Error(),
		}, err
	}
	if !success {
		return gen.Login401JSONResponse{
			Error: "invalid credentials",
		}, nil
	}

	id := email
	accessToken, err := h.ctrl.GenerateAccessTokenByID(id)
	if err != nil {
		return gen.Login500JSONResponse{
			Error: "failed to generate access token",
		}, err
	}
	accessTokenStr := string(accessToken)

	refreshToken, err := h.ctrl.GenerateRefreshToken()
	if err != nil {
		return gen.Login500JSONResponse{
			Error: "failed to generate refresh token",
		}, err
	}
	refreshTokenStr := string(refreshToken)
	refreshTokenMaxAge := h.ctrl.RefreshTokenMaxAge()
	refreshTokenEndPoint := fmt.Sprintf("/%s/%s", constant.APIResponseVersionV1, h.ctrl.RefreshTokenRefreshEndPoint())
	setCookie := fmt.Sprintf("refresh_token=%s; HttpOnly; Secure; SameSite=Lax; Path=%s; Max-Age=%d", refreshTokenStr, refreshTokenEndPoint, refreshTokenMaxAge)

	return gen.Login200JSONResponse{
		Body: gen.AuthResponse{
			AccessToken: &accessTokenStr,
		},
		Headers: gen.Login200ResponseHeaders{
			VersionId: constant.APIResponseVersionV1,
			SetCookie: setCookie,
		},
	}, nil
}

// Logout is the handler for the Logout endpoint.
func (h *Handler) Logout(_ context.Context, _ gen.LogoutRequestObject) (gen.LogoutResponseObject, error) {
	return gen.Logout204Response{
		Headers: gen.Logout204ResponseHeaders{
			VersionId: constant.APIResponseVersionV1,
		},
	}, nil
}
