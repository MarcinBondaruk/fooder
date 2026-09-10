package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/MarcinBondaruk/fooder/internal/ingredient"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) CreateRecipe(ctx context.Context, rcp recipe.Recipe) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		"INSERT INTO recipes (name, description) VALUES (?, ?)",
		rcp.Title, rcp.Description,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert recipe: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve recipe id: %w", err)
	}

	for _, ri := range rcp.Ingredients {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, unit) VALUES (?, ?, ?, ?)",
			id, ri.IngredientID, ri.Amount, ri.Unit.String(),
		)
		if err != nil {
			return 0, fmt.Errorf("failed to insert recipe ingredient: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(id), nil
}

func (r *RecipeRepository) GetRecipe(ctx context.Context, id int) (recipe.Recipe, error) {
	query := `
		SELECT r.id, r.name, r.description, ri.ingredient_id, i.name, ri.amount, ri.unit
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE r.id = ?`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("failed to query recipe: %w", err)
	}
	defer rows.Close()

	var rcp recipe.Recipe
	found := false

	for rows.Next() {
		var ingredientID *int
		var ingredientName *string
		var amount *float64
		var unit *string

		if err := rows.Scan(&rcp.ID, &rcp.Title, &rcp.Description, &ingredientID, &ingredientName, &amount, &unit); err != nil {
			return recipe.Recipe{}, fmt.Errorf("failed to scan recipe: %w", err)
		}
		found = true

		if ingredientID != nil {
			rcp.Ingredients = append(rcp.Ingredients, recipe.RecipeIngredient{
				IngredientID: *ingredientID,
				Name:         *ingredientName,
				Amount:       *amount,
				Unit:         ingredient.Unit(*unit),
			})
		}
	}

	if !found {
		return recipe.Recipe{}, recipe.ErrRecipeNotFound
	}

	return rcp, nil
}

func (r *RecipeRepository) GetRecipesByIds(ctx context.Context, ids []int) ([]recipe.Recipe, error) {
	if len(ids) == 0 {
		return []recipe.Recipe{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.name, r.description, ri.ingredient_id, i.name, ri.amount, ri.unit
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE r.id IN (%s)
		ORDER BY r.id`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

func (r *RecipeRepository) FindAllRecipes(ctx context.Context) ([]recipe.Recipe, error) {
	query := `
		SELECT r.id, r.name, r.description, ri.ingredient_id, i.name, ri.amount, ri.unit
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		ORDER BY r.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	return scanRecipes(rows)
}

func scanRecipes(rows *sql.Rows) ([]recipe.Recipe, error) {
	recipeMap := make(map[int]*recipe.Recipe)
	var order []int

	for rows.Next() {
		var recipeID int
		var name, description string
		var ingredientID *int
		var ingredientName *string
		var amount *float64
		var unit *string

		if err := rows.Scan(&recipeID, &name, &description, &ingredientID, &ingredientName, &amount, &unit); err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}

		rcp, exists := recipeMap[recipeID]
		if !exists {
			rcp = &recipe.Recipe{
				ID:          recipeID,
				Title:       name,
				Description: description,
			}
			recipeMap[recipeID] = rcp
			order = append(order, recipeID)
		}

		if ingredientID != nil {
			rcp.Ingredients = append(rcp.Ingredients, recipe.RecipeIngredient{
				IngredientID: *ingredientID,
				Name:         *ingredientName,
				Amount:       *amount,
				Unit:         ingredient.Unit(*unit),
			})
		}
	}

	recipes := make([]recipe.Recipe, 0, len(order))
	for _, id := range order {
		recipes = append(recipes, *recipeMap[id])
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
