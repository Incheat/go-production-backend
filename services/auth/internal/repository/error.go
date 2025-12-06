package repository

import "errors"

var (
	ErrRefreshTokenAlreadyExists = errors.New("refresh token already exists")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)