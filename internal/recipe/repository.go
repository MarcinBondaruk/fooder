package recipe

import "context"

type Repository interface {
	createRecipe(ctx context.Context, recipe Recipe) (int, error)
	getRecipe(ctx context.Context, id int) (Recipe, error)
	getRecipesByIds(ctx context.Context, ids []int) ([]Recipe, error)
	findAllRecipes(ctx context.Context) []Recipe
}
