package auth

import (
	"context"
	"errors"
)

type InMemoryTokenRepository struct {
	storage map[string]struct{}
}

func NewInMemoryTokenRepository(storage map[string]struct{}) *InMemoryTokenRepository {
	return &InMemoryTokenRepository{
		storage: storage,
	}
}

func (r *InMemoryTokenRepository) addToken(ctx context.Context, token string) error {
	r.storage[token] = struct{}{}

	return nil
}

func (r *InMemoryTokenRepository) findToken(ctx context.Context, token string) (string, error) {
	_, ok := r.storage[token]
	if !ok {
		return "", errors.New("no match")
	}

	return token, nil
}

func (r *InMemoryTokenRepository) deleteToken(ctx context.Context, token string) {
	delete(r.storage, token)
}
