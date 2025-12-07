// Package token defines the JWT maker for the auth service.
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/incheat/go-playground/services/auth/pkg/model"
)

// JWTMaker is a JWT maker.
type JWTMaker struct {
	secret []byte
	expire time.Duration
}

// NewJWTMaker creates a new JWT maker.
func NewJWTMaker(secret string, minutes int) *JWTMaker {
	return &JWTMaker{
		secret: []byte(secret),
		expire: time.Duration(minutes) * time.Minute,
	}
}

// CreateToken creates a new JWT token for a user.
func (m *JWTMaker) CreateToken(ID string) (model.AccessToken, error) {
	claims := jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(m.expire).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return model.AccessToken(token), nil
}

// ParseToken parses the ID from a JWT token.
// Returns the ID if the token is valid, otherwise returns an error.
func (m *JWTMaker) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKey
	}
	sub, _ := claims["sub"].(string)
	return sub, nil
}
