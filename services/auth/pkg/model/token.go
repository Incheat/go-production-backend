// Package model defines the models for the auth service.
package model

// AccessToken is a string that represents an access token.
type AccessToken string

// RefreshToken is a string that represents a refresh token.
type RefreshToken string

// type RefreshToken struct {
// 	ID        string `json:"id"`
// 	MemberID  string `json:"member_id"`
// 	TokenHash string `json:"token_hash"`
// 	ExpiresAt time.Time `json:"expires_at"`
// 	CreatedAt time.Time `json:"created_at"`
// 	RevokedAt time.Time `json:"revoked_at"`
// 	UserAgent string `json:"user_agent"`
// 	IPAddress string `json:"ip_address"`
// }
