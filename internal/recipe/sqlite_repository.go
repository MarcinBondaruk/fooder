package recipe

import (
	"database/sql"
	"fmt"
	"log"
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

func (r *SqliteRepository) CreateRecipe(name, description string, ingredients []string) int {
	stmt, err := r.db.Prepare("INSERT INTO recipes (name, description, ingredients) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	serializedIngredients := strings.Join(ingredients, ",")
	result, err := stmt.Exec(name, description, serializedIngredients)
	if err != nil {
		log.Fatalf("Failed to insert recipe: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to get last insert ID: %v", err)
	}

	return int(id)
}

func (r *SqliteRepository) FindRecipeById(id int) map[string]string {
	query := `SELECT id, name, description, ingredients FROM recipes WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var recipeID int
	var name, description, ingredients string
	err := row.Scan(&recipeID, &name, &description, &ingredients)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatalf("Failed to find recipe by ID: %v", err)
	}

	return map[string]string{
		"id":          fmt.Sprintf("%d", recipeID),
		"name":        name,
		"description": description,
		"ingredients": ingredients,
	}
}

func (r *SqliteRepository) GetRecipes(ids []int) []map[string]string {
	if len(ids) == 0 {
		return []map[string]string{}
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`SELECT id, name, description, ingredients FROM recipes WHERE id IN (%s)`, placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Fatalf("Failed to get recipes: %v", err)
	}
	defer rows.Close()

	var recipes []map[string]string
	for rows.Next() {
		var recipeID int
		var name, description, ingredients string
		err := rows.Scan(&recipeID, &name, &description, &ingredients)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		recipes = append(recipes, map[string]string{
			"id":          fmt.Sprintf("%d", recipeID),
			"name":        name,
			"description": description,
			"ingredients": ingredients,
		})
	}

	return recipes
}
