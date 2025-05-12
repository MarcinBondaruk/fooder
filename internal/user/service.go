package user

import (
	"context"
	"errors"
	"github.com/MarcinBondaruk/fooder/internal/auth"
)

type Service struct {
	authSvc    *auth.Service
	repository Repository
}

func NewService(authSvc *auth.Service, repository Repository) *Service {
	return &Service{
		authSvc:    authSvc,
		repository: repository,
	}
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repository.getUser(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = s.authSvc.Authenticate(ctx, password, user.Password)
	if err != nil {
		return "", err
	}

	token, err := s.authSvc.NewToken(ctx)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) LogoutUser(ctx context.Context, token string) {
	s.authSvc.DeleteToken(ctx, token)
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
