package recipe

type Repository interface {
	CreateRecipe(name, description string, ingredients []string) int
	FindRecipeById(id int) map[string]string
	GetRecipes(ids []int) []map[string]string
}
