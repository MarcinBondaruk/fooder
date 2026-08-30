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

func (r *InMemoryTokenRepository) AddToken(ctx context.Context, token string) error {
	r.storage[token] = struct{}{}

	return nil
}

func (r *InMemoryTokenRepository) FindToken(ctx context.Context, token string) (string, error) {
	_, ok := r.storage[token]
	if !ok {
		return "", errors.New("no match")
	}

	return token, nil
}

func (r *InMemoryTokenRepository) DeleteToken(ctx context.Context, token string) {
	delete(r.storage, token)
}
