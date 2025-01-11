package recipe

type Repository interface {
	createRecipe(recipe Recipe) (int, error)
	getRecipe(id int) (Recipe, error)
	getRecipesByIds(ids []int) ([]Recipe, error)
	findAllRecipes() []Recipe
}
