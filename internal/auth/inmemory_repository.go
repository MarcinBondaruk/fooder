package auth

import (
	"context"
	"errors"
)

type InMemoryRepository struct {
	storage map[string]struct{}
}

func NewInMemoryRepository(storage map[string]struct{}) *InMemoryRepository {
	return &InMemoryRepository{
		storage: storage,
	}
}

func (r *InMemoryRepository) addToken(ctx context.Context, token string) error {
	r.storage[token] = struct{}{}

	return nil
}

func (r *InMemoryRepository) findToken(ctx context.Context, token string) (string, error) {
	_, ok := r.storage[token]
	if !ok {
		return "", errors.New("no match")
	}

	return token, nil
}

func (r *InMemoryRepository) deleteToken(ctx context.Context, token string) {
	delete(r.storage, token)
}
