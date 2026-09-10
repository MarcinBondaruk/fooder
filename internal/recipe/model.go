package recipe

import "github.com/MarcinBondaruk/fooder/internal/ingredient"

type RecipeIngredient struct {
	IngredientID int
	Name         string
	Amount       float64
	Unit         ingredient.Unit
}

type Recipe struct {
	ID          int
	Title       string
	Description string
	Ingredients []RecipeIngredient
}

func NewRecipe(title string, description string, ingredients []RecipeIngredient) Recipe {
	return Recipe{
		Title:       title,
		Description: description,
		Ingredients: ingredients,
	}
}
