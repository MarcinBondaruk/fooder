package recipe

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository,
	}
}

func (s *Service) CreateRecipe(recipe Recipe) (int, error) {
	id, err := s.repository.createRecipe(recipe)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetRecipe(id int) (Recipe, error) {
	recipe, err := s.repository.getRecipe(id)

	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func (s *Service) GetRecipesByIds(ids []int) ([]Recipe, error) {
	recipes, err := s.repository.getRecipesByIds(ids)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (s *Service) FindAllRecipes() []Recipe {
	return s.repository.findAllRecipes()
}
