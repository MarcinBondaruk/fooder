package user

import "context"

type Repository interface {
	addUser(ctx context.Context, user User) (int, error)
	getUser(ctx context.Context, email string) (User, error)
}
