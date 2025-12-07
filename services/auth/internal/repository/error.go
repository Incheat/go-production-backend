// Package repository defines the errors for the auth service.
package repository

import "errors"

var (
	// ErrRefreshTokenAlreadyExists is the error for when a refresh token already exists.
	ErrRefreshTokenAlreadyExists = errors.New("refresh token already exists")
	// ErrRefreshTokenNotFound is the error for when a refresh token is not found.
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
