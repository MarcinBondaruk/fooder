package ingredient

import "context"

type Repository interface {
	CreateIngredient(ctx context.Context, ingredient Ingredient) (int, error)
	GetIngredient(ctx context.Context, id int) (Ingredient, error)
	FindAllIngredients(ctx context.Context) ([]Ingredient, error)
	FindOrCreate(ctx context.Context, name string) (Ingredient, error)
}
