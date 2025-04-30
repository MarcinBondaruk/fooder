package user

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/MarcinBondaruk/fooder/internal/auth"
)

type Service struct {
	authSvc *auth.Service
}

func NewService(authSvc *auth.Service) *Service {
	return &Service{
		authSvc: authSvc,
	}
}

func (s *Service) LoginUser(email, password string) (string, error) {
	if email != "admin@bendit.com" || password != "dupadupa" {
		return "", errors.New("login failed")
	}

	b := make([]byte, 32) // 256-bit
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	err := s.authSvc.StoreToken(token)
	if err != nil {
		return "", err
	}

	return token, nil
}
