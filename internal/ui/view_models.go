package ui

type HomeViewModel struct {
	Recipes []RecipeViewModel
}

type RecipeIngredientViewModel struct {
	Name   string
	Amount float64
	Unit   string
}

type RecipeViewModel struct {
	ID          int
	Name        string
	Description string
	Ingredients []RecipeIngredientViewModel
}
