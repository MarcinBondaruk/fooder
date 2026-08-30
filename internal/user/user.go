// Package user allows for manipulation of User entities
package user

type User struct {
	ID       int
	Email    string
	Name     *string
	Password string
}
