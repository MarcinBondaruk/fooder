package user

type Repository interface {
	addUser(user User) (int, error)
	getUser(email string) (User, error)
}
