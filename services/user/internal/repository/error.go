// Package repository defines the errors for the user service.
package repository

import "errors"

var (
	// ErrAlreadyExists is the error for when record already exists.
	ErrAlreadyExists = errors.New("record already exists")
	// ErrNotFound is the error for when record is not found.
	ErrNotFound = errors.New("record not found")
)
