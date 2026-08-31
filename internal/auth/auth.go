package auth

type UserCredentials struct {
	ID             int
	Email          string
	HashedPassword string
}
