// Package userhandler defines the server for the User GRPC API.
package userhandler

import (
	"context"

	userpb "github.com/incheat/go-production-backend/api/user/grpc/gen"
	userservice "github.com/incheat/go-production-backend/services/user/internal/service/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is the server for the User GRPC API.
type Server struct {
	service *userservice.Service
	userpb.UnimplementedUserServiceInternalServer
}

// New creates a new Server.
func New(service *userservice.Service) *Server {
	return &Server{service: service}
}

// VerifyUserCredentials is the server for the VerifyUserCredentials endpoint.
func (s *Server) VerifyUserCredentials(
	ctx context.Context,
	req *userpb.VerifyUserCredentialsRequest,
) (*userpb.VerifyUserCredentialsResponse, error) {

	email := req.Email
	password := req.Password

	user, err := s.service.VerifyUserCredentials(ctx, email, password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.VerifyUserCredentialsResponse{
		Id:     user.ID,
		Email:  user.Email,
		Status: user.Status,
	}, nil
}
