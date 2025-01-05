package main

import (
	"github.com/MarcinBondaruk/fooder/src/api"
	"github.com/MarcinBondaruk/fooder/src/cooking_list"
	"github.com/MarcinBondaruk/fooder/src/recipe"
	"log"
	"net/http"
)

func main() {
	recipeRepository := recipe.NewInMemoryRepository()
	cookingListRepository := cooking_list.NewInMemoryRepository()

	recipeSvc := recipe.NewService(recipeRepository)
	clSvc := cooking_list.NewService(cookingListRepository)

	http.HandleFunc("POST /api/v1/recipes", api.CreateRecipeHandler(recipeSvc))

	http.HandleFunc("GET /api/v1/recipes/{id}", api.ViewRecipeHandler(recipeSvc))

	http.HandleFunc("POST /api/v1/cooking-lists", api.CreateCookingListHandler(clSvc))

	http.HandleFunc("PATCH /api/v1/cooking-lists/{id}", api.AddRecipeToCookingListHandler(clSvc))

	http.HandleFunc("GET /api/v1/cooking-lists/{id}/shopping-list", api.GenerateShoppingListHandler(recipeSvc, clSvc))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
