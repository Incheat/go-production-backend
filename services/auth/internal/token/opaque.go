// Package token defines the Opaque maker for the auth service.
package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/incheat/go-production-backend/services/auth/pkg/model"
)

// OpaqueMaker is a Opaque maker.
type OpaqueMaker struct {
	numBytes        int
	maxAge          int
	refreshEndPoint string
}

// NewOpaqueMaker creates a new OpaqueMaker.
func NewOpaqueMaker(numBytes int, maxAge int, refreshEndPoint string) *OpaqueMaker {
	return &OpaqueMaker{numBytes: numBytes, maxAge: maxAge, refreshEndPoint: refreshEndPoint}
}

// CreateToken creates a URL-safe random token of given byte length.
func (m *OpaqueMaker) CreateToken() (model.RefreshToken, error) {
	b := make([]byte, m.numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}
	// URL-safe, no padding
	return model.RefreshToken(base64.RawURLEncoding.EncodeToString(b)), nil
}

// MaxAge returns the maximum age of the refresh token.
func (m *OpaqueMaker) MaxAge() int {
	return m.maxAge
}

// RefreshEndPoint returns the refresh end point.
func (m *OpaqueMaker) RefreshEndPoint() string {
	return m.refreshEndPoint
}
