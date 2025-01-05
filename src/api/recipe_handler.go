package api

import (
	"encoding/json"
	"fmt"
	"github.com/MarcinBondaruk/fooder/src/recipe"
	"net/http"
	"strconv"
	"strings"
)

func CreateRecipeHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRecipeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid contents", http.StatusBadRequest)
			return
		}

		id := recipeSvc.CreateRecipe(req.Name, req.Description, req.Ingredients)

		w.Header().Set("Location", fmt.Sprintf("/api/v1/recipes/%d", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func ViewRecipeHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		rcp, err := recipeSvc.GetRecipeById(id)
		if err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RecipeResponse{
			ID:          id,
			Name:        rcp["name"],
			Description: rcp["description"],
			Ingredients: strings.Split(rcp["ingredients"], ","),
		})
	}
}
