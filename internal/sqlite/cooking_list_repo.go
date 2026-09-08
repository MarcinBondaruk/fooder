package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
)

type CookingListRepository struct {
	db *sql.DB
}

func NewCookingListRepository(db *sql.DB) *CookingListRepository {
	query := `
		CREATE TABLE IF NOT EXISTS cooking_list (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    recipes TEXT NOT NULL
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("failed to initialize cooking list table", err)
	}

	return &CookingListRepository{
		db: db,
	}
}

func (r *CookingListRepository) CreateCookingList(ctx context.Context, cl cooking_list.CookingList) (int, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO main.cooking_list (recipes) VALUES (:recipes)`, serializeCookingListRecipes(cl.Recipes))
	if err != nil {
		return 0, errors.New("failed to insert cooking list")
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *CookingListRepository) UpdateCookingList(ctx context.Context, cl cooking_list.CookingList) error {
	serializedRecipes := serializeCookingListRecipes(cl.Recipes)

	_, err := r.db.ExecContext(ctx, "UPDATE main.cooking_list SET recipes = :recipes WHERE id = :id", serializedRecipes, cl.ID)
	if err != nil {
		return errors.New("failed to update a cooking list")
	}

	return nil
}

func (r *CookingListRepository) GetCookingList(ctx context.Context, cookingListID int) (cooking_list.CookingList, error) {
	var cl cooking_list.CookingList
	var serializedRecipes string
	result := r.db.QueryRowContext(ctx, "SELECT id, recipes FROM main.cooking_list WHERE id = :id", cookingListID)

	err := result.Scan(&cl.ID, &serializedRecipes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cooking_list.CookingList{}, errors.New("no cooking list found")
		}

		return cooking_list.CookingList{}, errors.New("failed to scan a cooking list")
	}

	cl.Recipes = deserializeCookingListRecipes(serializedRecipes)

	return cl, nil
}

func serializeCookingListRecipes(recipes []int) string {
	tmp := make([]string, len(recipes))
	for i, recipeID := range recipes {
		tmp[i] = strconv.Itoa(recipeID)
	}

	return strings.Join(tmp, ",")
}

func deserializeCookingListRecipes(serializedRecipes string) []int {
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
