package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type Ingredient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RecipeCreate struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
}

type RecipeResponse struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
}

func CreateRecipeHandler(logger *slog.Logger, recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RecipeCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Invalid request body",
			})
			return
		}

		ingredients := make([]string, len(req.Ingredients))
		for i, ing := range req.Ingredients {
			ingredients[i] = ing.Name
		}

		id, err := recipeSvc.CreateRecipe(r.Context(), recipe.NewRecipe(req.Title, req.Description, ingredients))
		if err != nil {
			logger.Error("failed to create recipe", "err", err)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/api/v1/recipes/%d", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(RecipeResponse{
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
			Ingredients: req.Ingredients,
		})
	}
}

func ViewRecipeHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Invalid recipe id",
			})
			return
		}

		rcp, err := recipeSvc.GetRecipe(r.Context(), id)
		if err != nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "Recipe not found",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toRecipeResponse(rcp))
	}
}

func ListRecipesHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rcps := recipeSvc.FindAllRecipes(r.Context())

		response := make([]RecipeResponse, len(rcps))
		for i, rcp := range rcps {
			response[i] = toRecipeResponse(rcp)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func toRecipeResponse(rcp recipe.Recipe) RecipeResponse {
	ingredients := make([]Ingredient, len(rcp.Ingredients))
	for i, name := range rcp.Ingredients {
		ingredients[i] = Ingredient{ID: i + 1, Name: name}
	}

	return RecipeResponse{
		ID:          rcp.ID,
		Title:       rcp.Title,
		Description: rcp.Description,
		Ingredients: ingredients,
	}
}
