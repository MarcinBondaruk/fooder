package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MarcinBondaruk/fooder/internal/ingredient"
)

type IngredientCreate struct {
	Name string `json:"name"`
}

type IngredientResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateIngredientHandler(logger *slog.Logger, ingredientSvc *ingredient.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IngredientCreate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeProblem(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		id, err := ingredientSvc.CreateIngredient(r.Context(), req.Name)
		if err != nil {
			logger.Error("failed to create ingredient", "err", err)
			writeProblem(w, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/api/v1/ingredients/%d", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(IngredientResponse{
			ID:   id,
			Name: req.Name,
		})
	}
}

func ViewIngredientHandler(logger *slog.Logger, ingredientSvc *ingredient.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeProblem(w, http.StatusBadRequest, "Invalid ingredient id")
			return
		}

		ing, err := ingredientSvc.GetIngredient(r.Context(), id)
		if err != nil {
			if errors.Is(err, ingredient.ErrIngredientNotFound) {
				writeProblem(w, http.StatusNotFound, "Ingredient not found")
				return
			}
			logger.Error("failed to get ingredient", "err", err)
			writeProblem(w, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IngredientResponse{
			ID:   ing.ID,
			Name: ing.Name,
		})
	}
}

func ListIngredientsHandler(logger *slog.Logger, ingredientSvc *ingredient.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := ingredientSvc.FindAllIngredients(r.Context())
		if err != nil {
			logger.Error("failed to list ingredients", "err", err)
			writeProblem(w, http.StatusInternalServerError, "")
			return
		}

		response := make([]IngredientResponse, len(ingredients))
		for i, ing := range ingredients {
			response[i] = IngredientResponse{ID: ing.ID, Name: ing.Name}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
