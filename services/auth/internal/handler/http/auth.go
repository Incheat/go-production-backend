// Package authhandler defines the server for the Auth API.
package authhandler

import (
	"context"
	"errors"
	"fmt"
	"log"

	servergen "github.com/incheat/go-production-backend/services/auth/internal/api/oapi/gen/public/server"
	"github.com/incheat/go-production-backend/services/auth/internal/constant"
	chimiddlewareutils "github.com/incheat/go-production-backend/services/auth/internal/middleware/chi/utils"
	authservice "github.com/incheat/go-production-backend/services/auth/internal/service/auth"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// _ is a placeholder to ensure that Server implements the StrictServerInterface interface.
var _ servergen.StrictServerInterface = (*Server)(nil)

// Server is the server for the Auth API.
type Server struct {
	service *authservice.Service
}

// New creates a new Server.
func New(service *authservice.Service) *Server {
	return &Server{service: service}
}

// Login is the server for the Login endpoint.
func (h *Server) Login(ctx context.Context, request servergen.LoginRequestObject) (servergen.LoginResponseObject, error) {

	sc := trace.SpanContextFromContext(ctx)
	log.Printf("[Login] trace_id=%s span_id=%s valid=%v",
		sc.TraceID().String(),
		sc.SpanID().String(),
		sc.IsValid(),
	)

	tr := otel.Tracer("auth.handler")
	ctx, span := tr.Start(ctx, "auth.login")
	defer span.End()

	email := string(request.Body.Email)
	password := request.Body.Password

	requestMeta, ok := chimiddlewareutils.GetRequestMeta(ctx)
	if !ok {
		return servergen.Login500JSONResponse{
			Error: "request metadata not found",
		}, errors.New("request metadata not found")
	}
	userAgent := requestMeta.UserAgent
	ipAddress := requestMeta.IPAddress

	res, err := h.service.LoginWithEmailAndPassword(ctx, email, password, userAgent, ipAddress)
	if err != nil {
		return servergen.Login500JSONResponse{
			Error: err.Error(),
		}, err
	}

	accessToken := string(res.AccessToken)
	setCookie := fmt.Sprintf("refresh_token=%s; HttpOnly; Secure; SameSite=Lax; Path=/%s/%s; Max-Age=%d", res.RefreshToken, constant.APIResponseVersionV1, res.RefreshEndPoint, res.RefreshMaxAgeSec)

	return servergen.Login200JSONResponse{
		Body: servergen.AuthResponse{
			AccessToken: &accessToken,
		},
		Headers: servergen.Login200ResponseHeaders{
			VersionId: constant.APIResponseVersionV1,
			SetCookie: setCookie,
		},
	}, nil
}

// Logout is the server for the Logout endpoint.
func (h *Server) Logout(_ context.Context, _ servergen.LogoutRequestObject) (servergen.LogoutResponseObject, error) {
	return servergen.Logout204Response{
		Headers: servergen.Logout204ResponseHeaders{
			VersionId: constant.APIResponseVersionV1,
		},
	}, nil
}
