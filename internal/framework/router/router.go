package router

import (
	"github.com/MarcinBondaruk/fooder/internal/api"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"github.com/MarcinBondaruk/fooder/internal/ui"
	"net/http"
)

func NewRouter(envs *env.Env, c *di.Container) *http.ServeMux {
	// middlewares
	commonMiddlewares := []middleware.Middleware{
		middleware.Logging,
		middleware.NewCorsMiddleware(envs.AllowedOrigins()),
	}

	commonAndAuthorizedMiddlewares := append(commonMiddlewares, middleware.NewApiKeyAuthorization(envs.ApiKey()))

	m := http.NewServeMux()

	// API
	// todo: add context support in handlers
	m.Handle("POST /api/v1/recipes", middleware.Chain(api.CreateRecipeHandler(c.RecipeService()), commonAndAuthorizedMiddlewares...))

	m.Handle("GET /api/v1/recipes/{id}", middleware.Chain(api.ViewRecipeHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("GET /api/v1/recipes", middleware.Chain(api.ListRecipesHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("POST /api/v1/cooking-lists", middleware.Chain(api.CreateCookingListHandler(c.CookingListService()), commonAndAuthorizedMiddlewares...))

	m.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(api.AddRecipeToCookingListHandler(c.CookingListService()), commonAndAuthorizedMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(api.ViewCookingListHandler(c.CookingListService()), commonMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(api.GenerateShoppingListHandler(c.CookingListService()), commonMiddlewares...))

	// PUBLIC UI
	m.Handle("GET /", http.RedirectHandler("/home", http.StatusFound))

	m.Handle("GET /home", middleware.Chain(ui.HomePageHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("GET /recipes/{id}", middleware.Chain(ui.RecipeDetailsPageHandler(c.RecipeService()), commonMiddlewares...))

	// ADMIN UI
	m.Handle("GET /admin/create-recipe", middleware.Chain(ui.ShowCreateRecipeForm(), commonMiddlewares...))

	m.Handle("POST /admin/recipes/create", middleware.Chain(ui.HandleCreateRecipe(c.RecipeService()), commonMiddlewares...))

	return m
}
