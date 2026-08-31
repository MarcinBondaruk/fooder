package router

import (
	"net/http"

	"github.com/MarcinBondaruk/fooder/internal/auth"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"github.com/MarcinBondaruk/fooder/internal/ui"
)

func NewRouter(envs *env.Env, c *di.Container) http.Handler {
	m := http.NewServeMux()

	// API
	m.Handle("POST /api/v1/recipes", middleware.Chain(recipe.CreateRecipeHandler(c.RecipeService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("GET /api/v1/recipes/{id}", middleware.Chain(recipe.ViewRecipeHandler(c.RecipeService())))

	m.Handle("GET /api/v1/recipes", middleware.Chain(recipe.ListRecipesHandler(c.RecipeService())))

	m.Handle("POST /api/v1/cooking-lists", middleware.Chain(cooking_list.CreateCookingListHandler(c.CookingListService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(cooking_list.AddRecipeToCookingListHandler(c.CookingListService()), middleware.NewApiKeyAuthorization(envs.ApiKey())))

	m.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(cooking_list.ViewCookingListHandler(c.CookingListService())))

	m.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(cooking_list.GenerateShoppingListHandler(c.CookingListService())))

	// PUBLIC UI
	m.Handle("GET /", http.RedirectHandler("/home", http.StatusFound))

	m.Handle("GET /favicon.ico", http.NotFoundHandler())

	m.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	m.Handle("GET /home", middleware.Chain(ui.HomePageHandler(c.RecipeService())))

	m.Handle("GET /recipes/{id}", middleware.Chain(ui.RecipeDetailsPageHandler(c.RecipeService())))

	// ADMIN UI
	m.Handle("GET /admin/login", middleware.Chain(ui.AdminLoginPage()))

	m.Handle("GET /admin/logout", middleware.Chain(auth.LogoutHandler(c.AuthService())))

	m.Handle("POST /admin/login-submit", middleware.Chain(auth.LoginSubmitHandler(c.LoginLimiter(), c.AuthService())))

	m.Handle("GET /admin/panel", middleware.Chain(ui.AdminPanelPage(), middleware.NewCookieBasedAuthorization(c.AuthService())))

	m.Handle("GET /admin/create-recipe", middleware.Chain(ui.ShowCreateRecipeForm(), middleware.NewCookieBasedAuthorization(c.AuthService())))

	m.Handle("POST /admin/recipes/create", middleware.Chain(ui.HandleCreateRecipe(c.RecipeService()), middleware.NewCookieBasedAuthorization(c.AuthService())))

	return middleware.Chain(m, middleware.PanicRecover(), middleware.LoggingMiddleware(), middleware.NewCorsMiddleware(envs.AllowedOrigins()))
}
