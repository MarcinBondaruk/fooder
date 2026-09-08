package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
)

type CreateCookingListRequest struct {
	RecipeID int `json:"recipeId"`
}

type AddRecipeToCookingListRequest struct {
	RecipeID int `json:"recipeId"`
}

type CookingListItem struct {
	RecipeID int `json:"recipeId"`
}

type CookingListResponse struct {
	ID      int               `json:"id"`
	Recipes []CookingListItem `json:"recipes"`
}

type ShoppingListResponse struct {
	ShoppingList map[string]int `json:"shoppingList"`
}

func CreateCookingListHandler(logger *slog.Logger, clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCookingListRequest
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

		id, err := clSvc.CreateCookingList(r.Context(), req.RecipeID)
		if err != nil {
			logger.Error("failed to create cooking list", "err", err)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/api/v1/cooking-lists/%d", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func ViewCookingListHandler(logger *slog.Logger, clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Invalid cooking list id",
			})
			return
		}

		cl, err := clSvc.ViewCookingList(r.Context(), id)
		if err != nil {
			if errors.Is(err, cooking_list.ErrCookingListNotFound) {
				w.Header().Set("Content-Type", "application/problem+json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ProblemJson{
					Type:   "about:blank",
					Title:  "Not Found",
					Status: http.StatusNotFound,
					Detail: "Cooking list not found",
				})
				return
			}
			logger.Error("failed to get cooking list", "err", err)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
			return
		}

		cookingListItems := make([]CookingListItem, len(cl.Recipes))
		for i, recipeId := range cl.Recipes {
			cookingListItems[i] = CookingListItem{
				RecipeID: recipeId,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CookingListResponse{
			ID:      cl.ID,
			Recipes: cookingListItems,
		})
	}
}

func AddRecipeToCookingListHandler(logger *slog.Logger, clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Invalid cooking list id",
			})
			return
		}

		var req AddRecipeToCookingListRequest
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

		err = clSvc.AddRecipeToCookingList(r.Context(), id, req.RecipeID)
		if err != nil {
			logger.Error("failed to add recipe to cooking list", "err", err)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Conflict",
				Status: http.StatusConflict,
				Detail: "Couldn't add recipe to cooking list",
			})
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func GenerateShoppingListHandler(logger *slog.Logger, clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Invalid cooking list id",
			})
			return
		}

		shoppingList, err := clSvc.GenerateShoppingList(r.Context(), id)
		if err != nil {
			logger.Error("failed to generate shopping list", "err", err)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProblemJson{
				Type:   "about:blank",
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ShoppingListResponse{
			ShoppingList: shoppingList,
		})
	}
}
