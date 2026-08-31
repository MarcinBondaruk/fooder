package auth

import (
	"context"
	"errors"
)

type TokenRepository interface {
	AddToken(ctx context.Context, token string) error
	FindToken(ctx context.Context, token string) (string, error)
	DeleteToken(ctx context.Context, token string)
}

type CredentialsRepository interface {
	GetCredentialsByEmail(ctx context.Context, email string) (UserCredentials, error)
}

var ErrCredentialsNotFound = errors.New("user credentials not found")
