package router

import (
	"github.com/MarcinBondaruk/fooder/internal/api"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"net/http"
)

func NewRouter(envs *env.Env, c *di.Container) *http.ServeMux {
	// middlewares
	commonMiddlewares := []middleware.Middleware{
		middleware.Logging,
		middleware.NewCorsMiddleware(envs.AllowedOrigins()),
	}

	m := http.NewServeMux()
	// todo: add context support in handlers
	m.Handle("POST /api/v1/recipes", middleware.Chain(api.CreateRecipeHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("GET /api/v1/recipes/{id}", middleware.Chain(api.ViewRecipeHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("GET /api/v1/recipes", middleware.Chain(api.ListRecipesHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("POST /api/v1/cooking-lists", middleware.Chain(api.CreateCookingListHandler(c.CookingListService()), commonMiddlewares...))

	m.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(api.AddRecipeToCookingListHandler(c.CookingListService()), commonMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(api.ViewCookingListHandler(c.CookingListService()), commonMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(api.GenerateShoppingListHandler(c.CookingListService()), commonMiddlewares...))

	return m
}
