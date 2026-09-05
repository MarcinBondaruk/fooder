package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type RecipeCreate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

type Recipe struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

func CreateRecipe(logger *slog.Logger, recipeSvc *recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload RecipeCreate
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			errResp := ProblemJson{
				Type:     "",
				Title:    "",
				Status:   http.StatusBadRequest,
				Detail:   "",
				Instance: "",
			}

			w.Header().Add("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		recipe, err := recipeSvc.CreateRecipe(r.Context(), recipe.NewRecipe(payload.Name, payload.Description, payload.Ingredients))
		if err != nil {
			errResp := ProblemJson{
				Type:     "",
				Title:    "",
				Status:   http.StatusBadRequest,
				Detail:   "",
				Instance: "",
			}

			w.Header().Add("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResp)
			return

		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipe)
	})
}
