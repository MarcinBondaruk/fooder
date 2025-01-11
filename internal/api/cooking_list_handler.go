package api

import (
	"encoding/json"
	"fmt"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"net/http"
	"strconv"
)

func CreateCookingListHandler(clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCookingListRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid contents", http.StatusBadRequest)
			return
		}

		id, err := clSvc.CreateCookingList(req.RecipeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/api/v1/cooking-lists/%d", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func ViewCookingListHandler(clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookingList, err := clSvc.ViewCookingList(id)
		cookingListItems := make([]CookingListItem, len(cookingList.Recipes()))
		for i, recipeId := range cookingList.Recipes() {
			cookingListItems[i] = CookingListItem{
				RecipeID: recipeId,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CookingListResponse{
			ID:      cookingList.ID,
			Recipes: cookingListItems,
		})
	}
}

func AddRecipeToCookingListHandler(clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		var req AddRecipeToCookingListRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid contents", http.StatusBadRequest)
			return
		}

		err = clSvc.AddRecipeToCookingList(id, req.RecipeID)

		if err != nil {
			http.Error(w, "Couldn't add recipe to cooking list", http.StatusConflict)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func GenerateShoppingListHandler(clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		shoppingList, err := clSvc.GenerateShoppingList(id)
		fmt.Printf("Shopping List: %+v\n", shoppingList)
		if err != nil {
			http.Error(w, "could not generate shopping list", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ShoppingListResponse{
			ShoppingList: shoppingList,
		})
	}
}
