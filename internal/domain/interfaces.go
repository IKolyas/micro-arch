package domain

type Repository[T any] interface {
	GetByID(id int) (*T, error)
	Find(fields map[string]string) ([]T, error)
	Save(*T) (res int, err error)
}

type UserRepository interface {
	Repository[User]
	GetAuthData(login string) (*User, error)
}
