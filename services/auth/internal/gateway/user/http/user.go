// Package usergateway defines the gateway for the user service.
package usergateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	clientgen "github.com/incheat/go-production-backend/api/user/oapi/gen/private"
	openapi_types "github.com/oapi-codegen/runtime/types"

	usermodel "github.com/incheat/go-production-backend/services/user/pkg/model"
)

// UserGateway is the gateway for the user service.
type UserGateway struct {
	client *clientgen.ClientWithResponses
}

// New creates a new UserGateway.
func New(addr string) (*UserGateway, error) {
	httpClient := &http.Client{Timeout: 2 * time.Second}

	c, err := clientgen.NewClientWithResponses(
		addr,
		clientgen.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}

	return &UserGateway{client: c}, nil
}

// VerifyCredentials verifies a user's credentials.
func (g *UserGateway) VerifyCredentials(ctx context.Context, email string, password string) (*usermodel.User, error) {
	// var editors []clientgen.RequestEditorFn
	// if token != "" {
	// 	editors = append(editors, func(ctx context.Context, req *http.Request) error {
	// 		req.Header.Set("Authorization", "Bearer "+token)
	// 		return nil
	// 	})
	// }

	resp, err := g.client.VerifyUserCredentialsWithResponse(
		ctx,
		clientgen.VerifyUserCredentialsJSONRequestBody{
			Email:    openapi_types.Email(email),
			Password: password,
		},
		// editors...,
	)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("user verify: 200 but empty JSON body")
		}
		return &usermodel.User{
			ID:     resp.JSON200.Id,
			Email:  string(resp.JSON200.Email),
			Status: resp.JSON200.Status,
		}, nil

	case http.StatusUnauthorized:
		if resp.JSON401 != nil {
			return nil, fmt.Errorf("user verify unauthorized: %s", resp.JSON401.Error)
		}
		return nil, fmt.Errorf("user verify unauthorized")

	default:
		return nil, fmt.Errorf("user verify unexpected status=%d body=%s", resp.StatusCode(), string(resp.Body))
	}
}
