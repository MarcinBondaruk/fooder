package auth

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) VerifyToken(token string) error {
	_, err := s.repository.findToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StoreToken(token string) error {
	return s.repository.addToken(token)
}
