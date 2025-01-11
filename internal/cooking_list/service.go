package cooking_list

import (
	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type Service struct {
	recipeSvc  *recipe.Service
	repository Repository
}

func NewService(recipeSvc *recipe.Service, repository Repository) *Service {
	return &Service{
		recipeSvc,
		repository,
	}
}

func (s *Service) CreateCookingList(recipeID int) (int, error) {
	cookingList := CookingList{recipes: []int{recipeID}}

	id, err := s.repository.createCookingList(cookingList)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) AddRecipeToCookingList(cookingListID, recipeID int) error {
	cookingList, err := s.repository.getCookingList(cookingListID)
	if err != nil {
		return err
	}

	cookingList.addRecipe(recipeID)

	err = s.repository.updateCookingList(cookingList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ViewCookingList(id int) (CookingList, error) {
	cookingList, err := s.repository.getCookingList(id)
	if err != nil {
		return CookingList{}, err
	}

	return cookingList, nil
}

func (s *Service) GenerateShoppingList(cookingListID int) (map[string]int, error) {
	cookingList, err := s.repository.getCookingList(cookingListID)
	if err != nil {
		return nil, err
	}

	recipes, err := s.recipeSvc.GetRecipesByIds(cookingList.Recipes())
	if err != nil {
		return nil, err
	}

	reducedIngredients := make(map[string]int)
	for _, r := range recipes {
		for _, ingredient := range r.Ingredients() {
			reducedIngredients[ingredient]++
		}
	}

	return reducedIngredients, nil
}
