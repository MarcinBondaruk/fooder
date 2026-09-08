package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	query := `
	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		ingredients TEXT NOT NULL
	)`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to create recipes table: %v", err)
	}

	return &RecipeRepository{
		db: db,
	}
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
		return 0, errors.New("failed to insert recipe into database: " + err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("failed to retrieve recipe id: " + err.Error())
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
			return recipe.Recipe{}, nil
		}
		log.Fatalf("Failed to find recipe by id: %v", err)
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
		log.Fatalf("Failed to get recipes: %v", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("error on rows close", "err", err)
		}
	}()

	var recipes []recipe.Recipe
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			return nil, errors.New("failed to scan recipes: " + err.Error())
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

func (r *RecipeRepository) FindAllRecipes(ctx context.Context) []recipe.Recipe {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, ingredients FROM main.recipes")
	if err != nil {
		log.Fatalf("Failed to get recipes: %v", err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("error on rows close", "err", err)
		}
	}()

	var recipes []recipe.Recipe
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		recipes = append(recipes, recipe.Recipe{
			ID:          recipeID,
			Title:       name,
			Description: description,
			Ingredients: strings.Split(ingredients, ","),
		})
	}

	return recipes
}

func serializeRecipeIds(recipes []int) string {
	tmp := make([]string, len(recipes))
	for i, recipeID := range recipes {
		tmp[i] = strconv.Itoa(recipeID)
	}

	return strings.Join(tmp, ",")
}
