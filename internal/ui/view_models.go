package ui

type HomeViewModel struct {
	Recipes []RecipeViewModel
}

type RecipeViewModel struct {
	ID          int
	Name        string
	Description string
	Ingredients []string
}
