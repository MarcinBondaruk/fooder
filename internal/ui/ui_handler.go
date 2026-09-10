package ui

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/MarcinBondaruk/fooder/internal/ingredient"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
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
			recipes[i] = toRecipeViewModel(rcp)
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

		viewModel := toRecipeViewModel(rcp)

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
		rawIngredients := strings.Split(r.FormValue("ingredients"), "\n")

		var ingredients []recipe.RecipeIngredient
		for _, line := range rawIngredients {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) < 3 {
				http.Error(w, fmt.Sprintf("invalid ingredient format: %s (expected: amount unit name)", line), http.StatusBadRequest)
				return
			}

			amount, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid amount: %s", parts[0]), http.StatusBadRequest)
				return
			}

			unit := ingredient.Unit(parts[1])
			if !unit.Valid() {
				http.Error(w, fmt.Sprintf("invalid unit: %s", parts[1]), http.StatusBadRequest)
				return
			}

			ingredientName := strings.Join(parts[2:], " ")
			ingredients = append(ingredients, recipe.RecipeIngredient{
				Name:   ingredientName,
				Amount: amount,
				Unit:   unit,
			})
		}

		_, err := recipeSvc.CreateRecipe(r.Context(), recipe.NewRecipe(name, description, ingredients))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

func toRecipeViewModel(rcp recipe.Recipe) RecipeViewModel {
	ingredients := make([]RecipeIngredientViewModel, len(rcp.Ingredients))
	for i, ri := range rcp.Ingredients {
		ingredients[i] = RecipeIngredientViewModel{
			Name:   ri.Name,
			Amount: ri.Amount,
			Unit:   ri.Unit.String(),
		}
	}
	return RecipeViewModel{
		ID:          rcp.ID,
		Name:        rcp.Title,
		Description: rcp.Description,
		Ingredients: ingredients,
	}
}
