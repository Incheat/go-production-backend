// Package model defines the models for the user service.
package model

import "time"

// Member is a model for a member.
type Member struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
