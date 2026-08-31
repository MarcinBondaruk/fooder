package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Service struct {
	credentialsRepo CredentialsRepository
	tokenRepo       TokenRepository
}

func NewService(credentialsRepo CredentialsRepository, tokenRepo TokenRepository) *Service {
	return &Service{
		credentialsRepo,
		tokenRepo,
	}
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {
	creds, err := s.credentialsRepo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if password != creds.HashedPassword {
		return "", errors.New("invalid credentials")
	}
	if err != nil {
		return "", err
	}

	b := make([]byte, 32) // 256-bit

	_, err = rand.Read(b)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	err = s.tokenRepo.AddToken(ctx, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) LogoutUser(ctx context.Context, token string) {
	s.tokenRepo.DeleteToken(ctx, token)
}

func (s *Service) VerifyToken(ctx context.Context, token string) error {
	_, err := s.tokenRepo.FindToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
