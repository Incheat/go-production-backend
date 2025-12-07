package token_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5" // If you're on v3/v4, adjust import path.
	"github.com/incheat/go-playground/services/auth/internal/token"
)

// TestCreateTokenAndParseID_RoundTrip checks the happy path: create a token and parse it back.
func TestCreateTokenAndParseID_RoundTrip(t *testing.T) {
	m := token.NewJWTMaker("test-secret", 15)
	userID := "user-123"

	token, err := m.CreateToken(userID)
	if err != nil {
		t.Fatalf("CreateToken error: %v", err)
	}
	if token == "" {
		t.Fatal("CreateToken returned empty token")
	}

	gotID, err := m.ParseToken(string(token))
	if err != nil {
		t.Fatalf("ParseID error: %v", err)
	}
	if gotID != userID {
		t.Fatalf("ParseID = %q, want %q", gotID, userID)
	}
}

// TestCreateToken_ExpClaimCloseToExpected checks that exp is roughly now + expire.
func TestCreateToken_ExpClaimCloseToExpected(t *testing.T) {
	const minutes = 15
	m := token.NewJWTMaker("test-secret", minutes)
	userID := "user-123"

	token, err := m.CreateToken(userID)
	if err != nil {
		t.Fatalf("CreateToken error: %v", err)
	}

	// Parse without verifying signature/claims to inspect exp.
	parsed, _, err := new(jwt.Parser).ParseUnverified(string(token), jwt.MapClaims{})
	if err != nil {
		t.Fatalf("ParseUnverified error: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims is not MapClaims")
	}

	expVal, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal(`exp claim missing or not a number`)
	}

	expTime := time.Unix(int64(expVal), 0)
	now := time.Now()

	// Allow a bit of slack in both directions.
	lowerBound := now.Add(time.Duration(minutes)*time.Minute - time.Minute)
	upperBound := now.Add(time.Duration(minutes)*time.Minute + time.Minute)

	if expTime.Before(lowerBound) || expTime.After(upperBound) {
		t.Fatalf("exp = %v not within [%v, %v]", expTime, lowerBound, upperBound)
	}
}

// TestParseID_InvalidToken ensures a totally invalid string fails.
func TestParseID_InvalidToken(t *testing.T) {
	m := token.NewJWTMaker("test-secret", 5)

	id, err := m.ParseToken("not-a-jwt")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
	if id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}
}

// TestParseID_WrongSecret ensures a token signed with another secret is rejected.
func TestParseID_WrongSecret(t *testing.T) {
	m1 := token.NewJWTMaker("secret-1", 5)
	m2 := token.NewJWTMaker("secret-2", 5)

	token, err := m1.CreateToken("user-123")
	if err != nil {
		t.Fatalf("CreateToken error: %v", err)
	}

	id, err := m2.ParseToken(string(token))
	if err == nil {
		t.Fatal("expected error when parsing token with wrong secret, got nil")
	}
	if id != "" {
		t.Fatalf("expected empty id, got %q", id)
	}
}

// TestParseID_MissingSubClaim documents behavior when "sub" is missing.
func TestParseID_MissingSubClaim(t *testing.T) {
	secret := "test-secret"
	m := token.NewJWTMaker(secret, 5)

	claims := jwt.MapClaims{
		// No "sub" on purpose
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := tkn.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("SignedString error: %v", err)
	}

	id, err := m.ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("ParseID error: %v", err)
	}
	if id != "" {
		t.Fatalf("expected empty id when sub is missing, got %q", id)
	}
}

// TestParseID_ExpiredToken ensures expired tokens are rejected.
func TestParseID_ExpiredToken(t *testing.T) {
	// Negative minutes => token's exp is already in the past.
	m := token.NewJWTMaker("test-secret", -1)

	token, err := m.CreateToken("user-123")
	if err != nil {
		t.Fatalf("CreateToken error: %v", err)
	}

	id, err := m.ParseToken(string(token))
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if id != "" {
		t.Fatalf("expected empty id for expired token, got %q", id)
	}
}
