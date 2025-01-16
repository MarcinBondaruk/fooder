package main

import (
	"database/sql"
	"github.com/MarcinBondaruk/fooder/internal/api"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	// init env wrapper
	envs, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	// todo: switch to postgres
	// todo: use env struct
	db, err := sql.Open("sqlite3", "/app/fooder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// todo: add DI
	recipeRepository := recipe.NewSqliteRepository(db)
	cookingListRepository := cooking_list.NewSqliteRepository(db)
	recipeSvc := recipe.NewService(recipeRepository)
	clSvc := cooking_list.NewService(recipeSvc, cookingListRepository)

	// middlewares
	commonMiddlewares := []middleware.Middleware{
		middleware.Logging,
		middleware.NewCorsMiddleware(envs.AllowedOrigins()),
	}

	// todo: encapsulate routes in router package
	// todo: add context support in handlers
	http.Handle("POST /api/v1/recipes", middleware.Chain(api.CreateRecipeHandler(recipeSvc), commonMiddlewares...))

	http.Handle("GET /api/v1/recipes/{id}", middleware.Chain(api.ViewRecipeHandler(recipeSvc), commonMiddlewares...))

	http.Handle("GET /api/v1/recipes", middleware.Chain(api.ListRecipesHandler(recipeSvc), commonMiddlewares...))

	http.Handle("POST /api/v1/cooking-lists", middleware.Chain(api.CreateCookingListHandler(clSvc), commonMiddlewares...))

	http.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(api.AddRecipeToCookingListHandler(clSvc), commonMiddlewares...))

	http.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(api.ViewCookingListHandler(clSvc), commonMiddlewares...))

	http.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(api.GenerateShoppingListHandler(clSvc), commonMiddlewares...))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
