package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MarcinBondaruk/fooder/internal/ingredient"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
)

type fakeRecipeRepo struct {
	createRecipeFn    func(ctx context.Context, r recipe.Recipe) (int, error)
	getRecipeFn       func(ctx context.Context, id int) (recipe.Recipe, error)
	getRecipesByIdsFn func(ctx context.Context, ids []int) ([]recipe.Recipe, error)
	findAllRecipesFn  func(ctx context.Context) ([]recipe.Recipe, error)
}

func (f *fakeRecipeRepo) CreateRecipe(ctx context.Context, r recipe.Recipe) (int, error) {
	return f.createRecipeFn(ctx, r)
}

func (f *fakeRecipeRepo) GetRecipe(ctx context.Context, id int) (recipe.Recipe, error) {
	return f.getRecipeFn(ctx, id)
}

func (f *fakeRecipeRepo) GetRecipesByIds(ctx context.Context, ids []int) ([]recipe.Recipe, error) {
	return f.getRecipesByIdsFn(ctx, ids)
}

func (f *fakeRecipeRepo) FindAllRecipes(ctx context.Context) ([]recipe.Recipe, error) {
	return f.findAllRecipesFn(ctx)
}

type fakeIngredientRepo struct {
	getIngredientFn func(ctx context.Context, id int) (ingredient.Ingredient, error)
}

func (f *fakeIngredientRepo) CreateIngredient(_ context.Context, _ ingredient.Ingredient) (int, error) {
	return 0, nil
}

func (f *fakeIngredientRepo) GetIngredient(ctx context.Context, id int) (ingredient.Ingredient, error) {
	return f.getIngredientFn(ctx, id)
}

func (f *fakeIngredientRepo) FindAllIngredients(_ context.Context) ([]ingredient.Ingredient, error) {
	return nil, nil
}

func (f *fakeIngredientRepo) FindOrCreate(_ context.Context, _ string) (ingredient.Ingredient, error) {
	return ingredient.Ingredient{}, nil
}

func newDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newFakeIngredientService() *ingredient.Service {
	return ingredient.NewService(&fakeIngredientRepo{
		getIngredientFn: func(_ context.Context, id int) (ingredient.Ingredient, error) {
			return ingredient.Ingredient{ID: id, Name: "test-ingredient"}, nil
		},
	})
}

func TestCreateRecipeHandler(t *testing.T) {
	tests := []struct {
		name            string
		body            string
		createFn        func(context.Context, recipe.Recipe) (int, error)
		wantStatus      int
		wantContentType string
		wantLocation    string
	}{
		{
			name: "valid recipe",
			body: `{"title":"Pasta","description":"Boil water","ingredients":[{"ingredientId":1,"amount":200,"unit":"g"}]}`,
			createFn: func(_ context.Context, _ recipe.Recipe) (int, error) {
				return 42, nil
			},
			wantStatus:      http.StatusCreated,
			wantContentType: "application/json",
			wantLocation:    "/api/v1/recipes/42",
		},
		{
			name:            "invalid json body",
			body:            `{invalid`,
			createFn:        nil,
			wantStatus:      http.StatusBadRequest,
			wantContentType: "application/problem+json",
		},
		{
			name:            "empty body",
			body:            "",
			createFn:        nil,
			wantStatus:      http.StatusBadRequest,
			wantContentType: "application/problem+json",
		},
		{
			name: "service error",
			body: `{"title":"Pasta","description":"Boil","ingredients":[]}`,
			createFn: func(_ context.Context, _ recipe.Recipe) (int, error) {
				return 0, errors.New("db down")
			},
			wantStatus:      http.StatusInternalServerError,
			wantContentType: "application/problem+json",
		},
		{
			name: "empty ingredients",
			body: `{"title":"Toast","description":"Simple","ingredients":[]}`,
			createFn: func(_ context.Context, _ recipe.Recipe) (int, error) {
				return 1, nil
			},
			wantStatus:      http.StatusCreated,
			wantContentType: "application/json",
			wantLocation:    "/api/v1/recipes/1",
		},
		{
			name:            "invalid unit",
			body:            `{"title":"Pasta","description":"Boil","ingredients":[{"ingredientId":1,"amount":200,"unit":"invalid"}]}`,
			createFn:        nil,
			wantStatus:      http.StatusBadRequest,
			wantContentType: "application/problem+json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRecipeRepo{createRecipeFn: tt.createFn}
			svc := recipe.NewService(repo)
			ingSvc := newFakeIngredientService()
			handler := CreateRecipeHandler(newDiscardLogger(), svc, ingSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/recipes", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			ct := rec.Header().Get("Content-Type")
			if ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			if tt.wantLocation != "" {
				loc := rec.Header().Get("Location")
				if loc != tt.wantLocation {
					t.Errorf("Location = %q, want %q", loc, tt.wantLocation)
				}
			}

			if tt.wantStatus == http.StatusCreated {
				var resp RecipeResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.Title == "" {
					t.Error("expected non-empty title in response")
				}
			}

			if tt.wantContentType == "application/problem+json" {
				var prob ProblemJson
				if err := json.NewDecoder(rec.Body).Decode(&prob); err != nil {
					t.Fatalf("failed to decode ProblemJson: %v", err)
				}
				if prob.Status != tt.wantStatus {
					t.Errorf("ProblemJson.Status = %d, want %d", prob.Status, tt.wantStatus)
				}
			}
		})
	}
}

func TestViewRecipeHandler(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string
		getRecipe  func(context.Context, int) (recipe.Recipe, error)
		wantStatus int
	}{
		{
			name:   "recipe found",
			pathID: "1",
			getRecipe: func(_ context.Context, _ int) (recipe.Recipe, error) {
				return recipe.Recipe{
					ID:          1,
					Title:       "Soup",
					Description: "Warm",
					Ingredients: []recipe.RecipeIngredient{
						{IngredientID: 1, Name: "water", Amount: 500, Unit: ingredient.UnitMilliliter},
						{IngredientID: 2, Name: "salt", Amount: 5, Unit: ingredient.UnitGram},
					},
				}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "recipe not found",
			pathID: "999",
			getRecipe: func(_ context.Context, _ int) (recipe.Recipe, error) {
				return recipe.Recipe{}, recipe.ErrRecipeNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id non-numeric",
			pathID:     "abc",
			getRecipe:  nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid id empty",
			pathID:     "",
			getRecipe:  nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRecipeRepo{getRecipeFn: tt.getRecipe}
			svc := recipe.NewService(repo)
			handler := ViewRecipeHandler(newDiscardLogger(), svc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/recipes/"+tt.pathID, nil)
			req.SetPathValue("id", tt.pathID)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp RecipeResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp.ID != 1 {
					t.Errorf("RecipeResponse.ID = %d, want 1", resp.ID)
				}
				if len(resp.Ingredients) != 2 {
					t.Errorf("len(Ingredients) = %d, want 2", len(resp.Ingredients))
				}
				if len(resp.Ingredients) > 0 && resp.Ingredients[0].IngredientID != 1 {
					t.Errorf("Ingredients[0].IngredientID = %d, want 1", resp.Ingredients[0].IngredientID)
				}
			}

			if tt.wantStatus == http.StatusBadRequest || tt.wantStatus == http.StatusNotFound {
				var prob ProblemJson
				if err := json.NewDecoder(rec.Body).Decode(&prob); err != nil {
					t.Fatalf("failed to decode ProblemJson: %v", err)
				}
				if prob.Status != tt.wantStatus {
					t.Errorf("ProblemJson.Status = %d, want %d", prob.Status, tt.wantStatus)
				}
			}
		})
	}
}

func TestListRecipesHandler(t *testing.T) {
	tests := []struct {
		name    string
		recipes []recipe.Recipe
		wantLen int
	}{
		{
			name: "multiple recipes",
			recipes: []recipe.Recipe{
				{ID: 1, Title: "Soup", Description: "Warm", Ingredients: []recipe.RecipeIngredient{
					{IngredientID: 1, Name: "water", Amount: 500, Unit: ingredient.UnitMilliliter},
				}},
				{ID: 2, Title: "Salad", Description: "Fresh", Ingredients: []recipe.RecipeIngredient{
					{IngredientID: 2, Name: "lettuce", Amount: 1, Unit: ingredient.UnitPiece},
					{IngredientID: 3, Name: "tomato", Amount: 2, Unit: ingredient.UnitPiece},
				}},
			},
			wantLen: 2,
		},
		{
			name:    "empty list",
			recipes: []recipe.Recipe{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRecipeRepo{
				findAllRecipesFn: func(_ context.Context) ([]recipe.Recipe, error) {
					return tt.recipes, nil
				},
			}
			svc := recipe.NewService(repo)
			handler := ListRecipesHandler(newDiscardLogger(), svc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/recipes", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
			}

			var resp []RecipeResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(resp) != tt.wantLen {
				t.Errorf("len(response) = %d, want %d", len(resp), tt.wantLen)
			}

			if tt.wantLen > 0 {
				second := resp[1]
				if len(second.Ingredients) != 2 {
					t.Errorf("second recipe ingredients len = %d, want 2", len(second.Ingredients))
				}
				if len(second.Ingredients) >= 2 && second.Ingredients[1].IngredientID != 3 {
					t.Errorf("second recipe Ingredients[1].IngredientID = %d, want 3", second.Ingredients[1].IngredientID)
				}
			}
		})
	}
}
