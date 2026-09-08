package recipe

import "context"

type RecipeRepository interface {
	CreateRecipe(ctx context.Context, recipe Recipe) (int, error)
	GetRecipe(ctx context.Context, id int) (Recipe, error)
	GetRecipesByIds(ctx context.Context, ids []int) ([]Recipe, error)
	FindAllRecipes(ctx context.Context) []Recipe
}
