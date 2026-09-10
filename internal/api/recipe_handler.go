package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MarcinBondaruk/fooder/internal/ingredient"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type RecipeIngredientRequest struct {
	IngredientID int     `json:"ingredientId"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
}

type RecipeIngredientResponse struct {
	IngredientID int     `json:"ingredientId"`
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
}

type RecipeCreateRequest struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Ingredients []RecipeIngredientRequest `json:"ingredients"`
}

type RecipeResponse struct {
	ID          int                        `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Ingredients []RecipeIngredientResponse `json:"ingredients"`
}

func CreateRecipeHandler(logger *slog.Logger, recipeSvc *recipe.Service, ingredientSvc *ingredient.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RecipeCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeProblem(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		recipeIngredients := make([]recipe.RecipeIngredient, len(req.Ingredients))
		for i, ri := range req.Ingredients {
			unit := ingredient.Unit(ri.Unit)
			if !unit.Valid() {
				writeProblem(w, http.StatusBadRequest, fmt.Sprintf("Invalid unit: %s", ri.Unit))
				return
			}

			ing, err := ingredientSvc.GetIngredient(r.Context(), ri.IngredientID)
			if err != nil {
				if errors.Is(err, ingredient.ErrIngredientNotFound) {
					writeProblem(w, http.StatusBadRequest, fmt.Sprintf("Ingredient with id %d not found", ri.IngredientID))
					return
				}
				logger.Error("failed to get ingredient", "err", err)
				writeProblem(w, http.StatusInternalServerError, "")
				return
			}

			recipeIngredients[i] = recipe.RecipeIngredient{
				IngredientID: ing.ID,
				Name:         ing.Name,
				Amount:       ri.Amount,
				Unit:         unit,
			}
		}

		id, err := recipeSvc.CreateRecipe(r.Context(), recipe.NewRecipe(req.Title, req.Description, recipeIngredients))
		if err != nil {
			logger.Error("failed to create recipe", "err", err)
			writeProblem(w, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/api/v1/recipes/%d", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(RecipeResponse{
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
			Ingredients: toRecipeIngredientResponses(recipeIngredients),
		})
	}
}

func ViewRecipeHandler(logger *slog.Logger, recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeProblem(w, http.StatusBadRequest, "Invalid recipe id")
			return
		}

		rcp, err := recipeSvc.GetRecipe(r.Context(), id)
		if err != nil {
			if errors.Is(err, recipe.ErrRecipeNotFound) {
				writeProblem(w, http.StatusNotFound, "Recipe not found")
				return
			}
			logger.Error("failed to get recipe", "err", err)
			writeProblem(w, http.StatusInternalServerError, "")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toRecipeResponse(rcp))
	}
}

func ListRecipesHandler(logger *slog.Logger, recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rcps, err := recipeSvc.FindAllRecipes(r.Context())
		if err != nil {
			logger.Error("failed to list recipes", "err", err)
			writeProblem(w, http.StatusInternalServerError, "")
			return
		}

		response := make([]RecipeResponse, len(rcps))
		for i, rcp := range rcps {
			response[i] = toRecipeResponse(rcp)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func toRecipeResponse(rcp recipe.Recipe) RecipeResponse {
	return RecipeResponse{
		ID:          rcp.ID,
		Title:       rcp.Title,
		Description: rcp.Description,
		Ingredients: toRecipeIngredientResponses(rcp.Ingredients),
	}
}

func toRecipeIngredientResponses(ingredients []recipe.RecipeIngredient) []RecipeIngredientResponse {
	responses := make([]RecipeIngredientResponse, len(ingredients))
	for i, ri := range ingredients {
		responses[i] = RecipeIngredientResponse{
			IngredientID: ri.IngredientID,
			Name:         ri.Name,
			Amount:       ri.Amount,
			Unit:         ri.Unit.String(),
		}
	}
	return responses
}
