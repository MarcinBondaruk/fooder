package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) NewToken() (string, error) {
	b := make([]byte, 32) // 256-bit

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	err = s.storeToken(token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) VerifyToken(token string) error {
	_, err := s.repository.findToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) storeToken(token string) error {
	return s.repository.addToken(token)
}

func (s *Service) Authenticate(password, userPassword string) error {
	if password != userPassword {
		return errors.New("invalid credentials")
	}

	return nil
}
