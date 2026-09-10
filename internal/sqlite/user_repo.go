package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MarcinBondaruk/fooder/internal/auth"
	"github.com/MarcinBondaruk/fooder/internal/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddUser(ctx context.Context, user user.User) (int, error) {
	result, err := r.db.ExecContext(
		ctx,
		"INSERT INTO main.users (email, name, password) VALUES (:email, :name, :password)",
		user.Email,
		user.Name,
		user.Password,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve user id: %w", err)
	}

	return int(id), nil
}

func (r *UserRepository) GetUser(ctx context.Context, email string) (*user.User, error) {
	u := user.User{}

	query := "SELECT id, email, name, password FROM main.users WHERE email = :email"
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) GetCredentialsByEmail(ctx context.Context, email string) (*auth.UserCredentials, error) {
	uc := auth.UserCredentials{}

	query := "SELECT id, email, password FROM main.users WHERE email = :email"
	row := ur.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&uc.ID, &uc.Email, &uc.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrCredentialsNotFound
		}

		return nil, err
	}

	return &uc, nil
}
