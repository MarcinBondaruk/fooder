package ingredient

import "context"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateIngredient(ctx context.Context, name string) (int, error) {
	return s.repository.CreateIngredient(ctx, NewIngredient(name))
}

func (s *Service) GetIngredient(ctx context.Context, id int) (Ingredient, error) {
	return s.repository.GetIngredient(ctx, id)
}

func (s *Service) FindAllIngredients(ctx context.Context) ([]Ingredient, error) {
	return s.repository.FindAllIngredients(ctx)
}

func (s *Service) FindOrCreate(ctx context.Context, name string) (Ingredient, error) {
	return s.repository.FindOrCreate(ctx, name)
}
