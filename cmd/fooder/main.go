package main

import (
	"database/sql"
	"github.com/MarcinBondaruk/fooder/internal/api"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "/app/fooder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// bootstrap, DI, Logging, whatever
	recipeRepository := recipe.NewSqliteRepository(db)
	cookingListRepository := cooking_list.NewSqliteRepository(db)
	recipeSvc := recipe.NewService(recipeRepository)
	clSvc := cooking_list.NewService(cookingListRepository)

	http.HandleFunc("POST /api/v1/recipes", api.CreateRecipeHandler(recipeSvc))

	http.HandleFunc("GET /api/v1/recipes/{id}", api.ViewRecipeHandler(recipeSvc))

	http.HandleFunc("GET /api/v1/recipes", api.ListRecipesHandler(recipeSvc))

	http.HandleFunc("POST /api/v1/cooking-lists", api.CreateCookingListHandler(clSvc))

	http.HandleFunc("PATCH /api/v1/cooking-lists/{id}", api.AddRecipeToCookingListHandler(clSvc))

	http.HandleFunc("GET /api/v1/cooking-lists/{id}/shopping-list", api.GenerateShoppingListHandler(recipeSvc, clSvc))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
