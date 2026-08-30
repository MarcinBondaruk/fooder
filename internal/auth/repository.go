package auth

import "context"

type TokenRepository interface {
	addToken(ctx context.Context, token string) error
	findToken(ctx context.Context, token string) (string, error)
	deleteToken(ctx context.Context, token string)
}

type CredentialsRepository interface {
	getUserCredentialsByEmail(ctx context.Context, email string) (UserCredentials, error)
}
