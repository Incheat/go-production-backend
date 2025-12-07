// Package handler defines the handlers for the Ping API.
package handler

import (
	"context"

	gen "github.com/incheat/go-playground/services/helloworld/internal/api/gen/server"
)

// _ is a placeholder to ensure that Server implements the StrictServerInterface interface.
var _ gen.StrictServerInterface = (*Server)(nil)

// Server is the server for the Ping API.
// It implements the StrictServerInterface interface.
type Server struct{}

// NewServer creates a new Server.
func NewServer() *Server {
	return &Server{}
}

// PingV1 is the handler for the PingV1 endpoint.
// It returns a 200 OK response with the message "pong" and the version ID "v1".
func (s *Server) PingV1(_ context.Context, _ gen.PingV1RequestObject) (gen.PingV1ResponseObject, error) {
	message := "pong"
	return gen.PingV1200JSONResponse{
		Body: gen.PingResponseV1{
			Message: &message,
		},
		Headers: gen.PingV1200ResponseHeaders{
			VersionId: "v1",
		},
	}, nil
}

// PingV2 is the handler for the PingV2 endpoint.
// It returns a 200 OK response with the message "pong" and the version ID "v2".
func (s *Server) PingV2(_ context.Context, _ gen.PingV2RequestObject) (gen.PingV2ResponseObject, error) {
	message := "pong"
	return gen.PingV2200JSONResponse{
		Body: gen.PingResponseV2{
			Message: &message,
		},
		Headers: gen.PingV2200ResponseHeaders{
			VersionId: "v2",
		},
	}, nil
}
