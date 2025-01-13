package api

import (
	"encoding/json"
	"fmt"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"net/http"
	"strconv"
)

func CreateRecipeHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRecipeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid contents", http.StatusBadRequest)
			return
		}

		id, err := recipeSvc.CreateRecipe(recipe.NewRecipe(req.Name, req.Description, req.Ingredients))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

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

		rcp, err := recipeSvc.GetRecipe(id)
		if err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RecipeResponse{
			ID:          id,
			Name:        rcp.Name(),
			Description: rcp.Description(),
			Ingredients: rcp.Ingredients(),
		})
	}
}

func ListRecipesHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rcps := recipeSvc.FindAllRecipes()

		responseRecipes := make([]RecipeResponse, len(rcps))

		for i, rcp := range rcps {
			responseRecipes[i] = RecipeResponse{
				ID:          rcp.ID(),
				Name:        rcp.Name(),
				Description: rcp.Description(),
				Ingredients: rcp.Ingredients(),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseRecipes)
	}
}
