package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/MarcinBondaruk/fooder/internal/framework/di"
	"github.com/MarcinBondaruk/fooder/internal/framework/env"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func main() {
	var (
		email    string
		password string
	)

	if len(os.Args) < 2 {
		fmt.Println("Expected command: create-user")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "create-user":
		createUserCmd := flag.NewFlagSet("create-user", flag.ExitOnError)
		e := createUserCmd.String("e", "", "User email")
		p := createUserCmd.String("p", "", "User password")

		createUserCmd.Parse(os.Args[2:])

		if *e == "" || *p == "" {
			fmt.Println("You must provide both email (-e) and password (-p)")
			createUserCmd.Usage()
			os.Exit(1)
		}

		email = *e
		password = *p

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Available commands: create-user")
		os.Exit(1)
	}

	// Use the parsed values
	envs, err := env.NewEnv()
	if err != nil {
		fmt.Println("Environment error:", err)
		os.Exit(1)
	}

	container, err := di.NewContainer(envs)
	if err != nil {
		fmt.Println("Failed to initialize container:", err)
		os.Exit(1)
	}

	id, err := container.UserService().CreateUser(context.Background(), email, password)
	if err != nil {
		fmt.Println("Failed to create user:", err)
		os.Exit(1)
	}

	fmt.Printf("%d", id)
}
