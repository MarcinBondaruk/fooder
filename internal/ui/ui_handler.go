package ui

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

func HomePageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		viewModel := HomeViewModel{
			Recipes: []RecipeViewModel{
				{
					ID:          1,
					Name:        "Woda z sola",
					Description: "Niedrogi studencki obiad",
					Ingredients: []string{"woda", "sól"},
				},
				{
					ID:          2,
					Name:        "Spaghetti Bolognese",
					Description: "Podstawa kuchni włoskiej",
					Ingredients: []string{"makaron spaghetti", "passata", "parmezan", "wołowe mięso mielone", "marchew", "cebula"},
				},
				{
					ID:          3,
					Name:        "Placki energetyczne",
					Description: "Niedrogi studencki obiad",
					Ingredients: []string{"mąka", "jajka", "czekolada", "banan", "serek wiejski", "proszek do pieczenia"},
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "home", viewModel)
	}
}

func RecipeDetailsPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		viewModel := struct {
			ID          int
			Name        string
			Description string
			Ingredients []string
		}{
			ID:          2,
			Name:        "Spaghetti Bolognese",
			Description: "Podstawa kuchni włoskiej",
			Ingredients: []string{"makaron spaghetti", "passata", "parmezan", "wołowe mięso mielone", "marchew", "cebula"},
		}

		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "recipe_details", viewModel)
	}
}
