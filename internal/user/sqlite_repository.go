package user

type SqliteRepository struct {
}

func NewSqliteRepository() *SqliteRepository {
	return &SqliteRepository{}
}

func (r *SqliteRepository) addUser(user User) error {
	return nil
}

func (r *SqliteRepository) getUser(email string) (User, error) {
	return User{}, nil
}
