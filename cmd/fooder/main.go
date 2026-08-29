package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/router"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGKILL, syscall.SIGTERM)

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

	go func() {
		log.Println("Starting http server on port 8080")
		err = s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}

		log.Println("server is shutting down")
	}()

	<-shutdownChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error during shutdown: %v", err)
	}

	log.Println("server shutdown gracefully")
}
