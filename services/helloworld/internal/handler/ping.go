// Package handler defines the handlers for the Ping API.
package handler

import (
	"context"

	"github.com/incheat/go-playground/services/helloworld/internal/api/gen/oapi/public/server"
)

// _ is a placeholder to ensure that Server implements the StrictServerInterface interface.
var _ server.StrictServerInterface = (*Server)(nil)

// Server is the server for the Ping API.
// It implements the StrictServerInterface interface.
type Server struct{}

// NewServer creates a new Server.
func NewServer() *Server {
	return &Server{}
}

// PingV1 is the handler for the PingV1 endpoint.
// It returns a 200 OK response with the message "pong" and the version ID "v1".
func (s *Server) PingV1(_ context.Context, _ server.PingV1RequestObject) (server.PingV1ResponseObject, error) {
	message := "pong"
	return server.PingV1200JSONResponse{
		Body: server.PingResponseV1{
			Message: &message,
		},
		Headers: server.PingV1200ResponseHeaders{
			VersionId: "v1",
		},
	}, nil
}

// PingV2 is the handler for the PingV2 endpoint.
// It returns a 200 OK response with the message "pong" and the version ID "v2".
func (s *Server) PingV2(_ context.Context, _ server.PingV2RequestObject) (server.PingV2ResponseObject, error) {
	message := "pong"
	return server.PingV2200JSONResponse{
		Body: server.PingResponseV2{
			Message: &message,
		},
		Headers: server.PingV2200ResponseHeaders{
			VersionId: "v2",
		},
	}, nil
}
