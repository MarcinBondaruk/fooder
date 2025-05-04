package user

import (
	"crypto/rand"
	"encoding/hex"
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

	b := make([]byte, 32) // 256-bit
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	err = s.authSvc.StoreToken(token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) CreateUser(email, password string) error {
	user := User{
		email,
		nil,
		password,
	}

	err := s.repository.addUser(user)
	if err != nil {
		return err
	}

	return nil
}
