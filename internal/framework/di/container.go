package di

import (
	"database/sql"
	"github.com/MarcinBondaruk/fooder/internal/auth"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"github.com/MarcinBondaruk/fooder/internal/user"
	"log"
)

type Services struct {
	authService    *auth.Service
	cookingService *cooking_list.Service
	recipeService  *recipe.Service
	userService    *user.Service
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

	tokenStorage := make(map[string]struct{})

	authRepository := auth.NewInMemoryRepository(tokenStorage)
	authSvc := auth.NewService(authRepository)

	userService := user.NewService(authSvc)

	recipeRepository := recipe.NewSqliteRepository(db)
	recipeSvc := recipe.NewService(recipeRepository)

	cookingListRepository := cooking_list.NewSqliteRepository(db)
	clSvc := cooking_list.NewService(recipeSvc, cookingListRepository)

	return &Container{
		db: db,
		services: &Services{
			authService:    authSvc,
			cookingService: clSvc,
			recipeService:  recipeSvc,
			userService:    userService,
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

func (c *Container) UserService() *user.Service {
	return c.services.userService
}

func (c *Container) AuthService() *auth.Service {
	return c.services.authService
}
