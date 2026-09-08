package recipe

import "context"

type Service struct {
	repository RecipeRepository
}

func NewService(repository RecipeRepository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, recipe Recipe) (int, error) {
	id, err := s.repository.CreateRecipe(ctx, recipe)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetRecipe(ctx context.Context, id int) (Recipe, error) {
	recipe, err := s.repository.GetRecipe(ctx, id)
	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func (s *Service) GetRecipesByIds(ctx context.Context, ids []int) ([]Recipe, error) {
	recipes, err := s.repository.GetRecipesByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (s *Service) FindAllRecipes(ctx context.Context) ([]Recipe, error) {
	return s.repository.FindAllRecipes(ctx)
}
