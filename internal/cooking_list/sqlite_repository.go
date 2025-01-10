package cooking_list

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
)

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) *SqliteRepository {
	cookingListSql := `
		CREATE TABLE IF NOT EXISTS cooking_list (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    recipes TEXT NOT NULL
		)
	`

	_, err := db.Exec(cookingListSql)
	if err != nil {
		log.Fatal("failed to initialize cooking list table", err)
	}

	return &SqliteRepository{
		db: db,
	}
}

func (r *SqliteRepository) createCookingList(cookingList CookingList) (int, error) {
	sql := `INSERT INTO cooking_list (recipes) VALUES (:recipes)`

	result, err := r.db.Exec(sql, serializeRecipes(cookingList.recipes))

	if err != nil {
		return 0, errors.New("failed to insert cooking list")
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *SqliteRepository) updateCookingList(cookingList CookingList) error {
	serializedRecipes := serializeRecipes(cookingList.recipes)

	_, err := r.db.Exec("UPDATE cooking_list SET recipes = :recipes WHERE id = :id", serializedRecipes, cookingList.ID)
	if err != nil {
		return errors.New("failed to update a cooking list")
	}

	return nil
}

func (r *SqliteRepository) getCookingList(cookingListID int) (CookingList, error) {
	var cookingList CookingList
	var serializedRecipes string
	result := r.db.QueryRow("SELECT id, recipes FROM cooking_list WHERE id = :id", cookingListID)

	err := result.Scan(&cookingList.ID, &serializedRecipes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CookingList{}, errors.New("no cooking list found")
		}

		return CookingList{}, errors.New("failed to scan a cooking list")
	}

	cookingList.recipes = deserializeRecipes(serializedRecipes)

	return cookingList, nil
}

func serializeRecipes(recipes []int) string {
	tmp := make([]string, len(recipes))
	for i, recipeID := range recipes {
		tmp[i] = strconv.Itoa(recipeID)
	}

	return strings.Join(tmp, ",")
}

func deserializeRecipes(serializedRecipes string) []int {
	tmp := strings.Split(serializedRecipes, ",")
	recipes := make([]int, len(tmp))

	for i, recipeID := range tmp {
		recipeID, err := strconv.Atoi(recipeID)
		if err != nil {
			log.Println("failed to deserialize recipe id:", recipeID, err)
		}
		recipes[i] = recipeID
	}

	return recipes
}
