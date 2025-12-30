// Package userhandler defines the server for the User API.
package userhandler

import (
	"context"

	servergen "github.com/incheat/go-production-backend/services/user/internal/api/oapi/gen/private/server"
	userservice "github.com/incheat/go-production-backend/services/user/internal/service/user"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// _ is a placeholder to ensure that Server implements the StrictServerInterface interface.
var _ servergen.StrictServerInterface = (*Server)(nil)

// Server is the server for the Auth API.
type Server struct {
	service *userservice.Service
}

// New creates a new Server.
func New(service *userservice.Service) *Server {
	return &Server{service: service}
}

// VerifyUserCredentials is the server for the VerifyUserCredentials endpoint.
func (s *Server) VerifyUserCredentials(ctx context.Context, request servergen.VerifyUserCredentialsRequestObject) (servergen.VerifyUserCredentialsResponseObject, error) {
	email := string(request.Body.Email)
	password := request.Body.Password

	user, err := s.service.VerifyUserCredentials(ctx, email, password)
	if err != nil {
		return servergen.VerifyUserCredentials401JSONResponse{
			Error: err.Error(),
		}, nil
	}

	return servergen.VerifyUserCredentials200JSONResponse{
		Id:     user.ID,
		Email:  openapi_types.Email(user.Email),
		Status: user.Status,
	}, nil
}
