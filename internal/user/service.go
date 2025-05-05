package user

import (
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

func (s *Service) LoginUser(email, password string) (string, error) {
	user, err := s.repository.getUser(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = s.authSvc.Authenticate(password, user.Password)
	if err != nil {
		return "", err
	}

	token, err := s.authSvc.NewToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) CreateUser(email, password string) error {
	user := User{
		Email:    email,
		Name:     nil,
		Password: password,
	}

	err := s.repository.addUser(user)
	if err != nil {
		return err
	}

	return nil
}
