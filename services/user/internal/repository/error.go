// Package repository defines the errors for the user service.
package repository

import "errors"

var (
	// ErrMemberAlreadyExists is the error for when a member already exists.
	ErrMemberAlreadyExists = errors.New("member already exists")
	// ErrMemberNotFound is the error for when a member is not found.
	ErrMemberNotFound = errors.New("member not found")
)
