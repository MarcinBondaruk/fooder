package auth

import (
	"context"
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

func (s *Service) NewToken(ctx context.Context) (string, error) {
	b := make([]byte, 32) // 256-bit

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	err = s.storeToken(ctx, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) VerifyToken(ctx context.Context, token string) error {
	_, err := s.repository.findToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) storeToken(ctx context.Context, token string) error {
	return s.repository.addToken(ctx, token)
}

func (s *Service) DeleteToken(ctx context.Context, token string) {
	s.repository.deleteToken(ctx, token)
}

func (s *Service) Authenticate(ctx context.Context, password, userPassword string) error {
	if password != userPassword {
		return errors.New("invalid credentials")
	}

	return nil
}
