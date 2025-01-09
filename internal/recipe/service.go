package recipe

import "errors"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) CreateRecipe(name, description string, ingredients []string) int {
	return s.repository.CreateRecipe(name, description, ingredients)
}

func (s *Service) GetRecipes(ids []int) []map[string]string {
	return s.repository.GetRecipes(ids)
}

func (s *Service) GetRecipeById(id int) (map[string]string, error) {
	recipe := s.repository.FindRecipeById(id)

	if recipe == nil {
		return nil, errors.New("recipe not found")
	}

	return recipe, nil
}
