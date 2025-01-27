package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
	"path"
	"time"
)

const (
	DefaultMigrationDir = "/db/migrations"
)

// Supports only PostgresSQL - for now
func main() {
	// check for migrations dir, fallback to default when necessary, create dir structure
	migrationsPath := os.Getenv("MIGRATIONS_DIR")
	if migrationsPath == "" {
		slog.Info("environment variable MIGRATIONS_DIR not set, will use default: ", "MIGRATIONS_DIR", DefaultMigrationDir)
		cwd, err := os.Getwd()
		if err != nil {
			slog.Error("could not determine current working directory", "error", err)
			return
		}

		migrationsPath = cwd + DefaultMigrationDir

		_, err = os.Stat(migrationsPath)
		if err != nil {
			slog.Info("migrations directory does not exists, create one")
			err = os.MkdirAll(migrationsPath, 0755)
			if err != nil {
				slog.Error("could not create migrations directory", "error", err)
				return
			}
		}
	}

	// create new migrations file
	timestamp := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS
	userFilename := os.Args[2]
	completeFilename := fmt.Sprintf("%s_%s.sql", timestamp, userFilename)
	newMigrationFile, err := os.Create(path.Join(migrationsPath, completeFilename))
	if err != nil {
		slog.Error("could not create migration file", "error", err)
		return
	}

	err = newMigrationFile.Close()
	if err != nil {
		slog.Error("could not close migration file", "error", err)
	}

	slog.Info("migration file created", "filename", completeFilename)

	// run migrations from given path
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		slog.Info("environment variable DB_HOST not set, will use default", "default", "localhost")
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		slog.Info("environment variable DB_PORT not set, will use default", "default", "5432")
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		slog.Info("environment variable DB_USER not set, will use default", "default", "postgres")
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		slog.Info("environment variable DB_PASSWORD not set, will use default", "default", "postgres")
		dbPassword = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		slog.Info("environment variable DB_NAME not set, will use default", "default", "postgres")
		dbName = "postgres"
	}

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		slog.Error("could not connect to database", "error", err)
		return
	}
	defer func() {
		err = conn.Close(context.Background())
		if err != nil {
			slog.Error("could not close connection", "error", err)
			return
		}
	}()

	// create migrations table
	createMigrationsTableSql := `
		CREATE TABLE IF NOT EXISTS migrations (
    		id SERIAL PRIMARY KEY,
    		name VARCHAR(255) NOT NULL UNIQUE,
		    executed_at TIMESTAMP NOT NULL
		);
	`

	_, err = conn.Exec(context.Background(), createMigrationsTableSql)
	if err != nil {
		slog.Error("could not initialize migrations table", "error", err)
		return
	}
}
