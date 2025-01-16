package main

import (
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/router"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	// init env wrapper
	log.Println("Initializing envs")
	envs, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	//todo: init dependency container
	//todo: init middlewares

	log.Println("Initializing router")
	r := router.NewRouter(envs)
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Starting http server on port 8080")
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal()
	}
}
