// Package usergateway defines the errors for the user gateway.
package usergateway

import "errors"

var (
	// ErrUserNotFound is the error for when the user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidPassword is the error for when the password is incorrect.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrUserBlocked is the error for when the user has been blocked (e.g. Service may need to handle this logic specially).
	ErrUserBlocked = errors.New("user is blocked")

	// ErrMalformedResponse is the error for remote service returned data in the wrong format (this is a technical error, but Gateway can annotate it).
	ErrMalformedResponse = errors.New(" malformed response")
)
