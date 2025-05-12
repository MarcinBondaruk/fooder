package router

import (
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/middleware"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"github.com/MarcinBondaruk/fooder/internal/ui"
	"github.com/MarcinBondaruk/fooder/internal/user"
	"net/http"
)

func NewRouter(envs *env.Env, c *di.Container) *http.ServeMux {
	// middlewares
	commonMiddlewares := []middleware.Middleware{
		middleware.Logging,
		middleware.NewCorsMiddleware(envs.AllowedOrigins()),
	}

	commonAndAuthorizedMiddlewares := append(commonMiddlewares, middleware.NewApiKeyAuthorization(envs.ApiKey()))
	commonAndUserAuth := append(commonMiddlewares, middleware.NewCookieBasedAuthorization(c.AuthService()))

	m := http.NewServeMux()

	// API
	m.Handle("POST /api/v1/recipes", middleware.Chain(recipe.CreateRecipeHandler(c.RecipeService()), commonAndAuthorizedMiddlewares...))

	m.Handle("GET /api/v1/recipes/{id}", middleware.Chain(recipe.ViewRecipeHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("GET /api/v1/recipes", middleware.Chain(recipe.ListRecipesHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("POST /api/v1/cooking-lists", middleware.Chain(cooking_list.CreateCookingListHandler(c.CookingListService()), commonAndAuthorizedMiddlewares...))

	m.Handle("PATCH /api/v1/cooking-lists/{id}", middleware.Chain(cooking_list.AddRecipeToCookingListHandler(c.CookingListService()), commonAndAuthorizedMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}", middleware.Chain(cooking_list.ViewCookingListHandler(c.CookingListService()), commonMiddlewares...))

	m.Handle("GET /api/v1/cooking-lists/{id}/shopping-list", middleware.Chain(cooking_list.GenerateShoppingListHandler(c.CookingListService()), commonMiddlewares...))

	// PUBLIC UI
	m.Handle("GET /", http.RedirectHandler("/home", http.StatusFound))

	m.Handle("GET /favicon.ico", http.NotFoundHandler())

	m.Handle("GET /home", middleware.Chain(ui.HomePageHandler(c.RecipeService()), commonMiddlewares...))

	m.Handle("GET /recipes/{id}", middleware.Chain(ui.RecipeDetailsPageHandler(c.RecipeService()), commonMiddlewares...))

	// ADMIN UI
	m.Handle("GET /admin/login", middleware.Chain(ui.AdminLoginPage(), commonMiddlewares...))

	m.Handle("GET /admin/logout", middleware.Chain(user.LogoutHandler(c.UserService()), commonMiddlewares...))

	m.Handle("POST /admin/login-submit", middleware.Chain(user.LoginSubmitHandler(c.LoginLimiter(), c.UserService()), commonMiddlewares...))

	m.Handle("GET /admin/panel", middleware.Chain(ui.AdminPanelPage(), commonAndUserAuth...))

	m.Handle("GET /admin/create-recipe", middleware.Chain(ui.ShowCreateRecipeForm(), commonAndUserAuth...))

	m.Handle("POST /admin/recipes/create", middleware.Chain(ui.HandleCreateRecipe(c.RecipeService()), commonAndUserAuth...))

	return m
}
