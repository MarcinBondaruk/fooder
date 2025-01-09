package api

import (
	"encoding/json"
	"fmt"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"net/http"
	"strconv"
	"strings"
)

func CreateCookingListHandler(clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := clSvc.CreateCookingList()
		w.Header().Set("Location", fmt.Sprintf("/api/v1/cooking-lists/%d", id))
		w.WriteHeader(http.StatusCreated)
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

func GenerateShoppingListHandler(recipeSvc *recipe.Service, clSvc *cooking_list.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		cookingList := clSvc.ViewCookingList(id)
		if len(cookingList) == 0 {
			json.NewEncoder(w).Encode(ShoppingListResponse{
				ShoppingList: []string{},
			})

			w.Header().Set("Content-Type", "application/json")
			return
		}

		recipes := recipeSvc.GetRecipes(cookingList)

		var shoppingList []string
		for _, rcp := range recipes {
			ingredients := strings.Split(rcp["ingredients"], ",")
			shoppingList = append(shoppingList, ingredients...)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ShoppingListResponse{
			ShoppingList: shoppingList,
		})
	}
}
