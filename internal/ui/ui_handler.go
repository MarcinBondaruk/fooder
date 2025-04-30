package ui

import (
	"embed"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

func HomePageHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rcps := recipeSvc.FindAllRecipes()
		recipes := make([]RecipeViewModel, len(rcps))

		for i, rcp := range rcps {
			recipes[i] = RecipeViewModel{
				ID:          rcp.ID(),
				Name:        rcp.Name(),
				Description: rcp.Description(),
				Ingredients: rcp.Ingredients(),
			}
		}

		viewModel := HomeViewModel{
			Recipes: recipes,
		}

		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "home", viewModel)
	}
}

func RecipeDetailsPageHandler(recipeSvc *recipe.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		rcp, err := recipeSvc.GetRecipe(id)
		if err != nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}

		viewModel := RecipeViewModel{
			ID:          id,
			Name:        rcp.Name(),
			Description: rcp.Description(),
			Ingredients: rcp.Ingredients(),
		}

		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "recipe_details", viewModel)
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

		_, err := recipeSvc.CreateRecipe(recipe.NewRecipe(name, description, ingredients))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
