package auth

import "errors"

type InMemoryRepository struct {
	storage map[string]struct{}
}

func NewInMemoryRepository(storage map[string]struct{}) *InMemoryRepository {
	return &InMemoryRepository{
		storage: storage,
	}
}

func (r *InMemoryRepository) addToken(token string) error {
	r.storage[token] = struct{}{}

	return nil
}

func (r *InMemoryRepository) findToken(token string) (string, error) {
	_, ok := r.storage[token]
	if !ok {
		return "", errors.New("no match")
	}

	return token, nil
}

func (r *InMemoryRepository) deleteToken(token string) {
	delete(r.storage, token)
}
