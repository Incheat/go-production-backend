// Package usergateway defines the gateway for the user service.
package usergateway

import (
	"context"
	"log"
	"net/http"
	"time"

	usermodel "github.com/incheat/go-playground/services/user/pkg/model"
)

// UserGateway is the gateway for the user service.
type UserGateway struct {
	baseURL    string
	httpClient *http.Client
}

// NewUserGateway creates a new UserGateway.
func NewUserGateway(baseURL string) *UserGateway {
	return &UserGateway{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// VerifyCredentials verifies a user's credentials.
func (g *UserGateway) VerifyCredentials(_ context.Context, _ string, _ string) (*usermodel.User, error) {
	req, err := http.NewRequest(http.MethodPost, g.baseURL+"/internal/users/verify", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close response body: %v", err)
		}
	}()

	var user usermodel.User
	return &user, nil
}
