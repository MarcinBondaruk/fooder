package cooking_list

import (
	"context"

	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type Service struct {
	recipeSvc  *recipe.Service
	repository CookingListRepository
}

func NewService(recipeSvc *recipe.Service, repository CookingListRepository) *Service {
	return &Service{
		recipeSvc,
		repository,
	}
}

func (s *Service) CreateCookingList(ctx context.Context, recipeID int) (int, error) {
	cookingList := NewCookingList([]int{recipeID})

	id, err := s.repository.CreateCookingList(ctx, cookingList)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) AddRecipeToCookingList(ctx context.Context, cookingListID, recipeID int) error {
	cookingList, err := s.repository.GetCookingList(ctx, cookingListID)
	if err != nil {
		return err
	}

	cookingList.AddRecipe(recipeID)

	err = s.repository.UpdateCookingList(ctx, cookingList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ViewCookingList(ctx context.Context, id int) (CookingList, error) {
	cookingList, err := s.repository.GetCookingList(ctx, id)
	if err != nil {
		return CookingList{}, err
	}

	return cookingList, nil
}

func (s *Service) GenerateShoppingList(ctx context.Context, cookingListID int) (map[string]int, error) {
	cookingList, err := s.repository.GetCookingList(ctx, cookingListID)
	if err != nil {
		return nil, err
	}

	recipes, err := s.recipeSvc.GetRecipesByIds(ctx, cookingList.Recipes)
	if err != nil {
		return nil, err
	}

	reducedIngredients := make(map[string]int)
	for _, r := range recipes {
		for _, ingredient := range r.Ingredients {
			reducedIngredients[ingredient]++
		}
	}

	return reducedIngredients, nil
}
