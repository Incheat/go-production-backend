// Package authservice defines the result for the auth API.
package authservice

import "github.com/incheat/go-production-backend/services/auth/pkg/model"

// LoginResult is the result for the login API.
type LoginResult struct {
	AccessToken      model.AccessToken
	RefreshToken     model.RefreshToken
	RefreshMaxAgeSec int
	RefreshEndPoint  string
	RefreshCookie    string
}
