package user

import (
	"database/sql"
	"errors"
	"log"
)

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) *SqliteRepository {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		name VARCHAR(63),
		password VARCHAR(63) NOT NULL
	)`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return &SqliteRepository{
		db: db,
	}
}

func (r *SqliteRepository) addUser(user User) (int, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (email, name, password) VALUES (:email, :name, :password)",
		user.Email,
		user.Name,
		user.Password,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("failed to retrieve user id" + err.Error())
	}

	return int(id), nil
}

func (r *SqliteRepository) getUser(email string) (User, error) {
	user := User{}

	query := "SELECT id, email, name, password FROM users WHERE email = :id"
	row := r.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		log.Fatalf("Failed to find user by email: %v", err)
	}

	return user, nil
}
