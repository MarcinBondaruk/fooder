package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/MarcinBondaruk/fooder/internal/auth/login_limiter"
)

type Service struct {
	credentialsRepo CredentialsRepository
	tokenRepo       TokenRepository
	loginLimiter    *login_limiter.LoginLimiter
}

func NewService(credentialsRepo CredentialsRepository, tokenRepo TokenRepository, loginLimiter *login_limiter.LoginLimiter) *Service {
	return &Service{
		credentialsRepo,
		tokenRepo,
		loginLimiter,
	}
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {
	creds, err := s.credentialsRepo.GetCredentialsByEmail(ctx, email)
	if err != nil {
		return "", ErrUserNotFound
	}

	if password != creds.HashedPassword {
		return "", ErrInvalidCredentials
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

func (s *Service) LoginUserWithIPLimit(ctx context.Context, ip, email, password string) (string, error) {
	t, err := s.LoginUser(ctx, email, password)
	if err != nil {
		if !s.loginLimiter.Register(ip) {
			return "", ErrTooManyAttempts
		}

		return "", err
	}

	s.loginLimiter.Release(ip)

	return t, nil
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
