package shopping_list

import "github.com/MarcinBondaruk/fooder/internal/ingredient"

type ShoppingListItem struct {
	IngredientID   int
	IngredientName string
	TotalAmount    float64
	Unit           ingredient.Unit
}
