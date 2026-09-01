package auth

import "errors"

var (
	ErrTooManyAttempts     = errors.New("too many login attempts")
	ErrCredentialsNotFound = errors.New("user credentials not found")
)
