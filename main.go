package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type CreateRecipeRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type RecipeResponse struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type InMemoryDB struct {
	lock    sync.Mutex
	recipes map[int]string
	nextID  int
}

func (db *InMemoryDB) CreateRecipe(contents string) int {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.nextID++
	db.recipes[db.nextID] = contents

	return db.nextID
}

func createRecipeHandler(db *InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRecipeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid contents", http.StatusBadRequest)
			return
		}

		id := db.CreateRecipe(req.Content)

		w.Header().Set("Location", fmt.Sprintf("/api/v1/recipes/%d", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func viewRecipeHandler(db *InMemoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		recipe, ok := db.recipes[id]
		if !ok {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RecipeResponse{ID: id, Content: recipe})
	}
}

func main() {
	db := &InMemoryDB{
		recipes: make(map[int]string),
	}

	http.HandleFunc("POST /api/v1/recipes", createRecipeHandler(db))

	http.HandleFunc("GET /api/v1/recipes/{id}", viewRecipeHandler(db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
