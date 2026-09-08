package ui

import (
	"embed"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

//go:embed templates/*.html templates/**/*.html
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html", "templates/**/*.html"))

func HomePageHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rcps, err := recipeSvc.FindAllRecipes(r.Context())
		if err != nil {
			slog.Error("failed to load recipes", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		recipes := make([]RecipeViewModel, len(rcps))

		for i, rcp := range rcps {
			recipes[i] = RecipeViewModel{
				ID:          rcp.ID,
				Name:        rcp.Title,
				Description: rcp.Description,
				Ingredients: rcp.Ingredients,
			}
		}

		viewModel := HomeViewModel{
			Recipes: recipes,
		}

		w.WriteHeader(http.StatusOK)
		err = templates.ExecuteTemplate(w, "home", viewModel)
		if err != nil {
			slog.Error("error parsing template", "error", err)
		}
	}
}

func RecipeDetailsPageHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		rcp, err := recipeSvc.GetRecipe(r.Context(), id)
		if err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		viewModel := RecipeViewModel{
			ID:          id,
			Name:        rcp.Title,
			Description: rcp.Description,
			Ingredients: rcp.Ingredients,
		}

		w.WriteHeader(http.StatusOK)
		err = templates.ExecuteTemplate(w, "recipe_details", viewModel)
		if err != nil {
			slog.Error("error parsing home template", "error", err)
		}
	}
}

func ShowCreateRecipeForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "admin_create_recipe", nil)
	}
}

func HandleCreateRecipe(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid data", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		description := r.FormValue("description")
		ingredients := strings.Split(r.FormValue("ingredients"), ",")

		_, err := recipeSvc.CreateRecipe(r.Context(), recipe.NewRecipe(name, description, ingredients))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
	}
}

func AdminLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "admin_login", nil)
		if err != nil {
			slog.Error("error parsing home template", "error", err)
		}
	}
}

func AdminPanelPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "admin_panel", nil)
		if err != nil {
			slog.Error("error parsing home template", "error", err)
		}
	}
}
