package recipe

type Recipe struct {
	ID          int
	Title       string
	Description string
	Ingredients []string
}

func NewRecipe(title string, description string, ingredients []string) Recipe {
	return Recipe{
		Title:       title,
		Description: description,
		Ingredients: ingredients,
	}
}
