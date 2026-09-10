package router

import (
	"net/http"

	"github.com/MarcinBondaruk/fooder/internal/api"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"github.com/MarcinBondaruk/fooder/internal/ui"
)

func NewRouter(envs *env.Env, c *di.Container) http.Handler {
	m := http.NewServeMux()

	// API
	m.Handle("POST /api/v1/ingredients", middleware.Chain(api.CreateIngredientHandler(c.Logger(), c.IngredientService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("GET /api/v1/ingredients/{id}", middleware.Chain(api.ViewIngredientHandler(c.Logger(), c.IngredientService())))

	m.Handle("GET /api/v1/ingredients", middleware.Chain(api.ListIngredientsHandler(c.Logger(), c.IngredientService())))

	m.Handle("POST /api/v1/recipes", middleware.Chain(api.CreateRecipeHandler(c.Logger(), c.RecipeService(), c.IngredientService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("GET /api/v1/recipes/{id}", middleware.Chain(api.ViewRecipeHandler(c.Logger(), c.RecipeService())))

	m.Handle("GET /api/v1/recipes", middleware.Chain(api.ListRecipesHandler(c.Logger(), c.RecipeService())))

	m.Handle("POST /api/v1/cooking-lists", middleware.Chain(api.CreateCookingListHandler(c.Logger(), c.CookingListService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(api.AddRecipeToCookingListHandler(c.Logger(), c.CookingListService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(api.ViewCookingListHandler(c.Logger(), c.CookingListService())))

	m.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(api.GenerateShoppingListHandler(c.Logger(), c.ShoppingListService())))

	m.Handle("GET /openapi.yaml", middleware.Chain(api.RawApiHandler()))

	m.Handle("GET /docs/api", middleware.Chain(api.SwaggerApiHandler()))

	// PUBLIC UI
	m.Handle("GET /", http.RedirectHandler("/home", http.StatusFound))

	m.Handle("GET /favicon.ico", http.NotFoundHandler())

	m.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	m.Handle("GET /home", middleware.Chain(ui.HomePageHandler(c.RecipeService())))

	m.Handle("GET /recipes/{id}", middleware.Chain(ui.RecipeDetailsPageHandler(c.RecipeService())))

	// ADMIN UI
	m.Handle("GET /admin/login", middleware.Chain(ui.AdminLoginPage()))

	m.Handle("GET /admin/logout", middleware.Chain(api.LogoutHandler(c.AuthService())))

	m.Handle("POST /admin/login-submit", middleware.Chain(api.LoginSubmitHandler(c.Logger(), c.AuthService())))

	m.Handle("GET /admin/panel", middleware.Chain(ui.AdminPanelPage(), middleware.NewCookieBasedAuthorization(c.AuthService())))

	m.Handle("GET /admin/create-recipe", middleware.Chain(ui.ShowCreateRecipeForm(), middleware.NewCookieBasedAuthorization(c.AuthService())))

	m.Handle("POST /admin/recipes/create", middleware.Chain(ui.HandleCreateRecipe(c.RecipeService()), middleware.NewCookieBasedAuthorization(c.AuthService())))

	return middleware.Chain(m, middleware.PanicRecover(), middleware.LoggingMiddleware(c.Logger()), middleware.NewCorsMiddleware(envs.AllowedOrigins()))
}
