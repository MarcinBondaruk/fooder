package main

import (
	"fmt"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	"github.com/MarcinBondaruk/fooder/internal/framework/router"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	// todo: bot protection - rate limiting, robots.txt

	// init env wrapper
	log.Println("Initializing envs")
	envs, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	c, err := di.NewContainer(envs)
	if err != nil {
		log.Fatal(err)
	}
	defer c.TearDown()

	//todo: init middlewares ?

	log.Println("Initializing router")
	r := router.NewRouter(envs, c)

	// todo: better server configuration max timeout and stuff
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", envs.ServerPort()),
		Handler: r,
	}

	log.Println("Starting http server on port 8080")
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal()
	}
}
