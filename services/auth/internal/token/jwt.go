package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/incheat/go-playground/services/auth/pkg/model"
)

// Maker is a JWT maker.
type Maker struct {
	secret []byte
	expire time.Duration
}

// NewMaker creates a new JWT maker.
func NewMaker(secret string, minutes int) *Maker {
	return &Maker{
		secret: []byte(secret),
		expire: time.Duration(minutes) * time.Minute,
	}
}

// CreateToken creates a new JWT token for a user.
func (m *Maker) CreateToken(ID string) (model.AccessToken, error) {
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

// ParseID parses the ID from a JWT token.
// Returns the ID if the token is valid, otherwise returns an error.
func (m *Maker) ParseID(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
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
