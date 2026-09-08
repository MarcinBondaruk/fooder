package auth

import "errors"

var (
	ErrTooManyAttempts     = errors.New("too many login attempts")
	ErrCredentialsNotFound = errors.New("user credentials not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
