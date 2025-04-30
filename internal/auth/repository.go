package auth

type Repository interface {
	addToken(token string) error
	findToken(token string) (string, error)
}
