package recipe

import (
	"log"
	"strconv"
	"strings"
	"sync"
)

type InMemoryRepository struct {
	lock    sync.Mutex
	recipes map[int]map[string]string
	nextID  int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		recipes: make(map[int]map[string]string),
	}
}

func (db *InMemoryRepository) CreateRecipe(name, description string, ingredients []string) int {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.nextID++
	newRecipe := make(map[string]string)
	newRecipe["id"] = strconv.Itoa(db.nextID)
	newRecipe["name"] = name
	newRecipe["description"] = description
	newRecipe["ingredients"] = strings.Join(ingredients, ",")

	db.recipes[db.nextID] = newRecipe

	return db.nextID
}

func (db *InMemoryRepository) FindRecipeById(id int) map[string]string {
	return db.recipes[id]
}

func (db *InMemoryRepository) FindRecipesByIds(ids []int) []map[string]string {
	db.lock.Lock()
	defer db.lock.Unlock()

	var recipes []map[string]string
	for _, recipeID := range ids {
		recipe, ok := db.recipes[recipeID]
		if !ok {
			log.Printf("Recipe %d not found", recipeID)
		}

		recipes = append(recipes, recipe)
	}

	return recipes
}
