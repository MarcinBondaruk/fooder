package user

type Repository interface {
	addUser(user User) error
	getUser(email string) (User, error)
}
