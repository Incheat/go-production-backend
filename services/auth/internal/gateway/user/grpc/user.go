// Package usergateway defines the gateway for the user service.
package usergateway

import (
	"context"
	"fmt"
	"time"

	userpb "github.com/incheat/go-production-backend/api/user/grpc/gen"
	usergateway "github.com/incheat/go-production-backend/services/auth/internal/gateway/user"
	usermodel "github.com/incheat/go-production-backend/services/user/pkg/model"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// UserGateway is the gateway for the user service.
type UserGateway struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceInternalClient
}

// New creates a new UserGateway.
func New(addr string) (*UserGateway, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &UserGateway{
		conn:   conn,
		client: userpb.NewUserServiceInternalClient(conn),
	}, nil
}

// Close closes the connection to the user service.
func (g *UserGateway) Close() error {
	return g.conn.Close()
}

// VerifyCredentials verifies a user's credentials.
func (g *UserGateway) VerifyCredentials(ctx context.Context, email string, password string) (usermodel.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := g.client.VerifyUserCredentials(ctx, &userpb.VerifyUserCredentialsRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return usermodel.User{}, g.translateError(err)
	}

	return usermodel.User{
		ID:     resp.GetId(),
		Email:  resp.GetEmail(),
		Status: resp.GetStatus(),
	}, nil
}

// translateError translates a gRPC error to an internal sentinel error.
func (g *UserGateway) translateError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		// If it's not a gRPC status, it might be a serious system error, wrap and throw it out.
		return fmt.Errorf("unexpected error type: %w", err)
	}

	switch st.Code() {
	// business errors: keep Sentinel, so that Service can make 401 judgment
	case codes.NotFound:
		return fmt.Errorf("%w: %v", usergateway.ErrUserNotFound, st.Message())
	case codes.Unauthenticated:
		return fmt.Errorf("%w: %v", usergateway.ErrInvalidPassword, st.Message())
	default:
		// For all other unexpected gRPC errors (e.g., Code.Internal, Code.Unavailable),
		// we wrap and return the original error.
		//
		// This preserves the error chain for cross-layer behavior detection without
		// tight coupling. For instance, if the downstream service is down or slow,
		// the returned error will satisfy errors.Is(err, context.DeadlineExceeded).
		//
		// By using %w, the Service layer can react to these standard library errors
		// (like timeouts) while remaining agnostic of the underlying gRPC transport.
		return fmt.Errorf("user service internal error [code=%v]: %w", st.Code(), err)
	}
}
