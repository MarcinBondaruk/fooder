package di

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/MarcinBondaruk/fooder/internal/auth"
	"github.com/MarcinBondaruk/fooder/internal/auth/login_limiter"
	"github.com/MarcinBondaruk/fooder/internal/cooking_list"
	"github.com/MarcinBondaruk/fooder/internal/framework/database"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/in_memory_db"
	"github.com/MarcinBondaruk/fooder/internal/recipe"
	"github.com/MarcinBondaruk/fooder/internal/sqlite"
	"github.com/MarcinBondaruk/fooder/internal/user"
)

type Services struct {
	authService    *auth.Service
	cookingService *cooking_list.Service
	recipeService  *recipe.Service
	userService    *user.Service
}

type Utils struct {
	logger       *slog.Logger
	loginLimiter *login_limiter.LoginLimiter
}

type Container struct {
	stopCh   chan struct{}
	db       *sql.DB
	services *Services
	utils    *Utils
}

func NewContainer(envs *env.Env) (*Container, error) {
	logHandlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
	logger := slog.New(logHandler)

	db, err := database.NewSqlite(envs.SqliteDSN())
	if err != nil {
		return nil, err
	}

	stopCh := make(chan struct{})
	loginLimiter := login_limiter.NewLoginLimiter(3, 5*time.Minute)
	loginLimiter.StartCleaner(stopCh)

	tokenStorage := make(map[string]struct{})

	tokenRepository := in_memory_db.NewTokenRepository(tokenStorage)
	userRepository := sqlite.NewUserRepository(db)

	authSvc := auth.NewService(userRepository, tokenRepository, loginLimiter)

	userService := user.NewService(userRepository)

	recipeRepository := sqlite.NewRecipeRepository(db)
	recipeSvc := recipe.NewService(recipeRepository)

	cookingListRepository := cooking_list.NewSqliteRepository(db)
	clSvc := cooking_list.NewService(recipeSvc, cookingListRepository)

	return &Container{
		stopCh: stopCh,
		db:     db,
		services: &Services{
			authService:    authSvc,
			cookingService: clSvc,
			recipeService:  recipeSvc,
			userService:    userService,
		},
		utils: &Utils{
			logger:       logger,
			loginLimiter: loginLimiter,
		},
	}, nil
}

func (c *Container) TearDown() error {
	err := c.db.Close()
	if err != nil {
		return err
	}

	close(c.stopCh)
	return nil
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

func (c *Container) Logger() *slog.Logger {
	return c.utils.logger
}

func (c *Container) LoginLimiter() *login_limiter.LoginLimiter {
	return c.utils.loginLimiter
}
