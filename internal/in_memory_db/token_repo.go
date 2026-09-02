package in_memory_db

import (
	"context"
	"errors"
)

type TokenRepository struct {
	storage map[string]struct{}
}

func NewTokenRepository(storage map[string]struct{}) *TokenRepository {
	return &TokenRepository{
		storage: storage,
	}
}

func (r *TokenRepository) AddToken(ctx context.Context, token string) error {
	r.storage[token] = struct{}{}

	return nil
}

func (r *TokenRepository) FindToken(ctx context.Context, token string) (string, error) {
	_, ok := r.storage[token]
	if !ok {
		return "", errors.New("no match")
	}

	return token, nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, token string) {
	delete(r.storage, token)
}
