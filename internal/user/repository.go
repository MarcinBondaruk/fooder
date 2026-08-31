package user

import (
	"context"
	"errors"
)

type Repository interface {
	AddUser(ctx context.Context, user User) (int, error)
	GetUser(ctx context.Context, email string) (*User, error)
}

var ErrUserNotFound = errors.New("user not found")
