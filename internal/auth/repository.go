package auth

import "context"

type Repository interface {
	addToken(ctx context.Context, token string) error
	findToken(ctx context.Context, token string) (string, error)
	deleteToken(ctx context.Context, token string)
}
