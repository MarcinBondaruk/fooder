package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
)

type CookingListRepository struct {
	db *sql.DB
}

func NewCookingListRepository(db *sql.DB) *CookingListRepository {
	return &CookingListRepository{db: db}
}

func (r *CookingListRepository) CreateCookingList(ctx context.Context, cl cooking_list.CookingList) (int, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO main.cooking_list (recipes) VALUES (:recipes)`, serializeCookingListRecipes(cl.Recipes))
	if err != nil {
		return 0, fmt.Errorf("failed to insert cooking list: %w", err)
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func (r *CookingListRepository) UpdateCookingList(ctx context.Context, cl cooking_list.CookingList) error {
	serializedRecipes := serializeCookingListRecipes(cl.Recipes)

	_, err := r.db.ExecContext(ctx, "UPDATE main.cooking_list SET recipes = :recipes WHERE id = :id", serializedRecipes, cl.ID)
	if err != nil {
		return fmt.Errorf("failed to update cooking list: %w", err)
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
			return cooking_list.CookingList{}, cooking_list.ErrCookingListNotFound
		}

		return cooking_list.CookingList{}, fmt.Errorf("failed to scan cooking list: %w", err)
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
	var recipes []int

	for _, s := range tmp {
		recipeID, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		recipes = append(recipes, recipeID)
	}

	return recipes
}
