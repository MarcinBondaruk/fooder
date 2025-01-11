package recipe

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) *SqliteRepository {
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

	return &SqliteRepository{
		db: db,
	}
}

func (r *SqliteRepository) createRecipe(recipe Recipe) (int, error) {
	serializedIngredients := strings.Join(recipe.ingredients, ",")
	result, err := r.db.Exec(
		"INSERT INTO recipes (name, description, ingredients) VALUES (:name, :description, :ingredients)",
		recipe.name,
		recipe.description,
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

func (r *SqliteRepository) getRecipe(id int) (Recipe, error) {
	query := `SELECT id, name, description, ingredients FROM recipes WHERE id = :id`
	row := r.db.QueryRow(query, id)

	var recipeID int
	var name, description, ingredients string
	err := row.Scan(&recipeID, &name, &description, &ingredients)
	if err != nil {
		if err == sql.ErrNoRows {
			return Recipe{}, nil
		}
		log.Fatalf("Failed to find recipe by id: %v", err)
	}

	return Recipe{
		recipeID,
		name,
		description,
		strings.Split(ingredients, ","),
	}, nil
}

func (r *SqliteRepository) getRecipesByIds(ids []int) ([]Recipe, error) {
	if len(ids) == 0 {
		return []Recipe{}, nil
	}

	query := fmt.Sprintf("SELECT id, name, description, ingredients FROM recipes WHERE id IN (%s)", serializeRecipeIds(ids))
	rows, err := r.db.Query(query)
	if err != nil {
		log.Fatalf("Failed to get recipes: %v", err)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			return nil, errors.New("failed to scan recipes: " + err.Error())
		}

		fmt.Printf("populating: %+v\n", recipeID)

		recipes = append(recipes, Recipe{
			recipeID,
			name,
			description,
			strings.Split(ingredients, ","),
		})
	}

	return recipes, nil
}

func (r *SqliteRepository) findAllRecipes() []Recipe {
	rows, err := r.db.Query("SELECT id, name, description, ingredients FROM recipes")
	if err != nil {
		log.Fatalf("Failed to get recipes: %v", err)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		recipes = append(recipes, Recipe{
			recipeID,
			name,
			description,
			strings.Split(ingredients, ","),
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
