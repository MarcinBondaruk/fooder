package shopping_list

import (
	"context"

	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/ingredient"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type CookingListService interface {
	ViewCookingList(ctx context.Context, id int) (cooking_list.CookingList, error)
}

type RecipeService interface {
	GetRecipesByIds(ctx context.Context, ids []int) ([]recipe.Recipe, error)
}

type Service struct {
	cookingListSvc CookingListService
	recipeSvc      RecipeService
}

func NewService(cookingListSvc CookingListService, recipeSvc RecipeService) *Service {
	return &Service{
		cookingListSvc: cookingListSvc,
		recipeSvc:      recipeSvc,
	}
}

func (s *Service) GenerateShoppingList(ctx context.Context, cookingListID int) ([]ShoppingListItem, error) {
	cl, err := s.cookingListSvc.ViewCookingList(ctx, cookingListID)
	if err != nil {
		return nil, err
	}

	recipes, err := s.recipeSvc.GetRecipesByIds(ctx, cl.Recipes)
	if err != nil {
		return nil, err
	}

	type aggKey struct {
		ingredientID int
		unit         ingredient.Unit
	}

	aggregated := make(map[aggKey]*ShoppingListItem)

	for _, r := range recipes {
		for _, ri := range r.Ingredients {
			key := aggKey{ingredientID: ri.IngredientID, unit: ri.Unit}
			if item, exists := aggregated[key]; exists {
				item.TotalAmount += ri.Amount
			} else {
				aggregated[key] = &ShoppingListItem{
					IngredientID:   ri.IngredientID,
					IngredientName: ri.Name,
					TotalAmount:    ri.Amount,
					Unit:           ri.Unit,
				}
			}
		}
	}

	result := make([]ShoppingListItem, 0, len(aggregated))
	for _, item := range aggregated {
		result = append(result, *item)
	}

	return result, nil
}
