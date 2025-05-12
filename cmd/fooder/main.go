package main

import (
	"fmt"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/router"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logHandlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// init env wrapper
	log.Println("Initializing envs...")
	envs, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	c, err := di.NewContainer(envs)
	if err != nil {
		log.Fatal(err)
	}
	defer c.TearDown()

	log.Println("Initializing router...")
	r := router.NewRouter(envs, c)

	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", envs.ServerPort()),
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	log.Println("Starting http server on port 8080")
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal()
	}
}
