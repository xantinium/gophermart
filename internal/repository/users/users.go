package usersrepo

type UsersRepositoryOptions struct {
	Storage UsersStorage
}

// NewUsersRepository создаёт новый репозиторий метрик.
func NewUsersRepository(opts UsersRepositoryOptions) *UsersRepository {
	return &UsersRepository{
		storage: opts.Storage,
	}
}

// UsersRepository структура, описывающая репозиторий метрик.
type UsersRepository struct {
	storage UsersStorage
}
