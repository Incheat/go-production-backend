// Package repository defines the errors for the user service.
package repository

import "errors"

var (
	// ErrUserAlreadyExists is the error for when a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound is the error for when a user is not found.
	ErrUserNotFound = errors.New("user not found")
)
