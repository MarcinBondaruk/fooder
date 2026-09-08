package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) (*RecipeRepository, error) {
	query := `
	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		ingredients TEXT NOT NULL
	)`
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create recipes table: %w", err)
	}

	return &RecipeRepository{
		db: db,
	}, nil
}

func (r *RecipeRepository) CreateRecipe(ctx context.Context, rcp recipe.Recipe) (int, error) {
	serializedIngredients := strings.Join(rcp.Ingredients, ",")
	result, err := r.db.ExecContext(
		ctx,
		"INSERT INTO main.recipes (name, description, ingredients) VALUES (:name, :description, :ingredients)",
		rcp.Title,
		rcp.Description,
		serializedIngredients,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert recipe: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve recipe id: %w", err)
	}

	return int(id), nil
}

func (r *RecipeRepository) GetRecipe(ctx context.Context, id int) (recipe.Recipe, error) {
	query := `SELECT id, name, description, ingredients FROM main.recipes WHERE id = :id`
	row := r.db.QueryRowContext(ctx, query, id)

	var recipeID int
	var name, description, ingredients string
	err := row.Scan(&recipeID, &name, &description, &ingredients)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recipe.Recipe{}, recipe.ErrRecipeNotFound
		}
		return recipe.Recipe{}, fmt.Errorf("failed to get recipe: %w", err)
	}

	return recipe.Recipe{
		ID:          recipeID,
		Title:       name,
		Description: description,
		Ingredients: strings.Split(ingredients, ","),
	}, nil
}

func (r *RecipeRepository) GetRecipesByIds(ctx context.Context, ids []int) ([]recipe.Recipe, error) {
	if len(ids) == 0 {
		return []recipe.Recipe{}, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name, description, ingredients FROM main.recipes WHERE id IN (:recipeIds)",
		serializeRecipeIds(ids),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []recipe.Recipe
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}

		recipes = append(recipes, recipe.Recipe{
			ID:          recipeID,
			Title:       name,
			Description: description,
			Ingredients: strings.Split(ingredients, ","),
		})
	}

	return recipes, nil
}

func (r *RecipeRepository) FindAllRecipes(ctx context.Context) ([]recipe.Recipe, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, ingredients FROM main.recipes")
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	var recipes []recipe.Recipe
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}

		recipes = append(recipes, recipe.Recipe{
			ID:          recipeID,
			Title:       name,
			Description: description,
			Ingredients: strings.Split(ingredients, ","),
		})
	}

	return recipes, nil
}

func serializeRecipeIds(recipes []int) string {
	tmp := make([]string, len(recipes))
	for i, recipeID := range recipes {
		tmp[i] = strconv.Itoa(recipeID)
	}

	return strings.Join(tmp, ",")
}
