package user

import (
	"context"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (int, error) {
	user := User{
		Email:    email,
		Name:     nil,
		Password: password,
	}

	id, err := s.repository.addUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
