package recipe

type Repository interface {
	CreateRecipe(name, description string, ingredients []string) int
	FindRecipeById(id int) map[string]string
	FindRecipesByIds(ids []int) []map[string]string
	FindAllRecipes() []map[string]string
}
