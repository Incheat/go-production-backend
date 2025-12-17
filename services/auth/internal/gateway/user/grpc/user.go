// Package usergateway defines the gateway for the user service.
package usergateway

import (
	"context"
	"time"

	userpb "github.com/incheat/go-playground/api/user/grpc/gen"
	usermodel "github.com/incheat/go-playground/services/user/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
func (g *UserGateway) VerifyCredentials(ctx context.Context, email string, password string) (*usermodel.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := g.client.VerifyUserCredentials(ctx, &userpb.VerifyUserCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &usermodel.User{
		ID:     resp.GetId(),
		Email:  resp.GetEmail(),
		Status: resp.GetStatus(),
	}, nil
}
