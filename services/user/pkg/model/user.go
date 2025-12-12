// Package model defines the models for the user service.
package model

import "time"

// User is a model for a user.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
