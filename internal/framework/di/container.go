package di

import (
	"database/sql"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"log"
)

type Services struct {
	cookingService *cooking_list.Service
	recipeService  *recipe.Service
}

type Container struct {
	db       *sql.DB
	services *Services
}

func NewContainer(envs *env.Env) (*Container, error) {
	db, err := sql.Open("sqlite3", envs.SqliteDSN())
	if err != nil {
		return nil, err
	}

	recipeRepository := recipe.NewSqliteRepository(db)
	cookingListRepository := cooking_list.NewSqliteRepository(db)
	recipeSvc := recipe.NewService(recipeRepository)
	clSvc := cooking_list.NewService(recipeSvc, cookingListRepository)

	return &Container{
		db: db,
		services: &Services{
			cookingService: clSvc,
			recipeService:  recipeSvc,
		},
	}, nil
}

func (c *Container) TearDown() {
	err := c.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Container) RecipeService() *recipe.Service {
	return c.services.recipeService
}

func (c *Container) CookingListService() *cooking_list.Service {
	return c.services.cookingService
}
