package user

import "context"

type Repository interface {
	AddUser(ctx context.Context, user User) (int, error)
	GetUser(ctx context.Context, email string) (User, error)
}
