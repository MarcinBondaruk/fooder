package router

import (
	"database/sql"
	"github.com/MarcinBondaruk/fooder/internal/api"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"log"
	"net/http"
)

func NewRouter(envs *env.Env) *http.ServeMux {
	// todo: switch to postgres
	// todo: use env struct
	db, err := sql.Open("sqlite3", "/app/fooder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// todo: add DI, build DI in main.go
	recipeRepository := recipe.NewSqliteRepository(db)
	cookingListRepository := cooking_list.NewSqliteRepository(db)
	recipeSvc := recipe.NewService(recipeRepository)
	clSvc := cooking_list.NewService(recipeSvc, cookingListRepository)

	// middlewares
	commonMiddlewares := []middleware.Middleware{
		middleware.Logging,
		middleware.NewCorsMiddleware(envs.AllowedOrigins()),
	}

	m := http.NewServeMux()
	// todo: add context support in handlers
	m.Handle("POST /api/v1/recipes", middleware.Chain(api.CreateRecipeHandler(recipeSvc), commonMiddlewares...))

	m.Handle("GET /api/v1/recipes/{id}", middleware.Chain(api.ViewRecipeHandler(recipeSvc), commonMiddlewares...))

	m.Handle("GET /api/v1/recipes", middleware.Chain(api.ListRecipesHandler(recipeSvc), commonMiddlewares...))

	m.Handle("POST /api/v1/cooking-lists", middleware.Chain(api.CreateCookingListHandler(clSvc), commonMiddlewares...))

	m.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(api.AddRecipeToCookingListHandler(clSvc), commonMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(api.ViewCookingListHandler(clSvc), commonMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(api.GenerateShoppingListHandler(clSvc), commonMiddlewares...))

	return m
}
