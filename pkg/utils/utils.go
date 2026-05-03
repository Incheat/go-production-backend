// Package utils defines the utilities for the project.
package utils

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}
