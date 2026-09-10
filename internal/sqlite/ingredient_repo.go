package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MarcinBondaruk/fooder/internal/ingredient"
)

type IngredientRepository struct {
	db *sql.DB
}

func NewIngredientRepository(db *sql.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

func (r *IngredientRepository) CreateIngredient(ctx context.Context, ing ingredient.Ingredient) (int, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO ingredients (name) VALUES (?)", ing.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to insert ingredient: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve ingredient id: %w", err)
	}

	return int(id), nil
}

func (r *IngredientRepository) GetIngredient(ctx context.Context, id int) (ingredient.Ingredient, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM ingredients WHERE id = ?", id)

	var ing ingredient.Ingredient
	err := row.Scan(&ing.ID, &ing.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ingredient.Ingredient{}, ingredient.ErrIngredientNotFound
		}
		return ingredient.Ingredient{}, fmt.Errorf("failed to get ingredient: %w", err)
	}

	return ing, nil
}

func (r *IngredientRepository) FindAllIngredients(ctx context.Context) ([]ingredient.Ingredient, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM ingredients")
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []ingredient.Ingredient
	for rows.Next() {
		var ing ingredient.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient: %w", err)
		}
		ingredients = append(ingredients, ing)
	}

	return ingredients, nil
}

func (r *IngredientRepository) FindOrCreate(ctx context.Context, name string) (ingredient.Ingredient, error) {
	_, err := r.db.ExecContext(ctx, "INSERT OR IGNORE INTO ingredients (name) VALUES (?)", name)
	if err != nil {
		return ingredient.Ingredient{}, fmt.Errorf("failed to insert ingredient: %w", err)
	}

	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM ingredients WHERE name = ?", name)

	var ing ingredient.Ingredient
	if err := row.Scan(&ing.ID, &ing.Name); err != nil {
		return ingredient.Ingredient{}, fmt.Errorf("failed to get ingredient: %w", err)
	}

	return ing, nil
}
