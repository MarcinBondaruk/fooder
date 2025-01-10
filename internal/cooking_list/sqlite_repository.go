package cooking_list

import (
	"database/sql"
	"fmt"
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

func (r *SqliteRepository) CreateCookingList(recipeID int) int {
	sql := `INSERT INTO cooking_list (recipes) VALUES (:recipes)`

	result, err := r.db.Exec(sql, strconv.Itoa(recipeID))

	if err != nil {
		log.Println("failed to insert cooking list", err)
		return 0
	}

	id, _ := result.LastInsertId()
	return int(id)
}

func (r *SqliteRepository) AddRecipeToCookingList(cookingListID, recipeID int) {
	queryResult, err := r.db.Query("SELECT recipes FROM cooking_list WHERE id = :id", cookingListID)
	if err != nil {
		log.Println("failed to retrieve cooking list", err)
		return
	}
	defer queryResult.Close()

	var recipesResult string
	if queryResult.Next() {
		err := queryResult.Scan(&recipesResult)
		if err != nil {
			log.Println("failed to scan a cooking list", err)
			return
		}
	}
	queryResult.Close()

	recipeStringIds := strings.Split(recipesResult, ",")
	recipeStringIds = append(recipeStringIds, strconv.Itoa(recipeID))

	_, err = r.db.Exec("UPDATE cooking_list SET recipes = :recipes WHERE id = :id", strings.Join(recipeStringIds, ","), cookingListID)
	if err != nil {
		log.Println("failed to update a cooking list", err)
		return
	}
}

func (r *SqliteRepository) ViewCookingList(cookingListID int) []int {
	var recipesResult string
	result := r.db.QueryRow("SELECT recipes FROM cooking_list WHERE id = :id", cookingListID)

	err := result.Scan(&recipesResult)
	if err != nil {
		log.Println("failed to scan a cooking list", err)
	}

	recipesStringIds := strings.Split(recipesResult, ",")
	fmt.Println(recipesStringIds)
	recipeIDs := make([]int, len(recipesStringIds))

	for i, id := range recipesStringIds {
		intID, err := strconv.Atoi(id)
		if err != nil {
			log.Println("failed to convert string to int", err)
		}
		recipeIDs[i] = intID
	}

	return recipeIDs
}
