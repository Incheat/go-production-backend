package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/incheat/go-production-backend/services/auth/pkg/model"
)

// JWTMaker is a JWT maker.
type JWTMaker struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	kid        string
	issuer     string
	audience   string
	expire     time.Duration
}

// New creates a new JWTMaker from a PEM encoded private key.
func New(privateKeyPEM, kid, issuer, audience string, expire time.Duration) (*JWTMaker, error) {
	if privateKeyPEM == "" {
		return nil, errors.New("JWT privateKeyPEM is empty")
	}
	if kid == "" {
		return nil, errors.New("JWT kid is empty")
	}
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM")
	}

	var priv *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv = key
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		priv, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("PKCS8 key is not RSA")
		}
	default:
		return nil, errors.New("unsupported PEM block type: " + block.Type)
	}

	return &JWTMaker{
		privateKey: priv,
		publicKey:  &priv.PublicKey,
		kid:        kid,
		issuer:     issuer,
		audience:   audience,
		expire:     expire,
	}, nil
}

// CreateToken creates a new RS256 JWT token for a user.
func (m *JWTMaker) CreateToken(ID string) (model.AccessToken, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": ID,
		"iss": m.issuer,
		"aud": m.audience, // ["user-api", "order-api", "auth-api"]
		"iat": now.Unix(),
		"exp": now.Add(m.expire).Unix(),
		// "scope": "user:read order:read auth:read"
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = m.kid
	tokenStr, err := t.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt string with private key: %w", err)
	}
	accessToken := model.AccessToken(tokenStr)
	return accessToken, nil
}

// ---- JWKS ----

type jwks struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Kty string `json:"kty"` // "RSA"
	Use string `json:"use"` // "sig"
	Alg string `json:"alg"` // "RS256"
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKSJSON returns the JWKS JSON for the public key.
func (m *JWTMaker) JWKSJSON() ([]byte, error) {
	j := jwks{
		Keys: []jwkKey{rsaPublicKeyToJWK(m.publicKey, m.kid)},
	}
	return json.Marshal(j)
}

// JWKSHandler returns the JWKS JSON for the public key.
func (m *JWTMaker) JWKSHandler(w http.ResponseWriter, _ *http.Request) {
	b, err := m.JWKSJSON()
	if err != nil {
		http.Error(w, "failed to build jwks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func rsaPublicKeyToJWK(pub *rsa.PublicKey, kid string) jwkKey {
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(intToBytes(pub.E))
	return jwkKey{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: kid,
		N:   n,
		E:   e,
	}
}

func intToBytes(i int) []byte {
	if i == 0 {
		return []byte{0}
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte(i & 0xff)}, b...)
		i >>= 8
	}
	return b
}
