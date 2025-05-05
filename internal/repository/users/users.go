package usersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

func New(storage UsersStorage) *UsersRepository {
	return &UsersRepository{
		storage: storage,
	}
}

type UsersRepository struct {
	storage UsersStorage
}

func (repo *UsersRepository) CreateUser(ctx context.Context, login, passwordHash string) error {
	return repo.storage.CreateUser(ctx, login, passwordHash)
}

func (repo *UsersRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	return repo.storage.GetUserByLogin(ctx, login)
}
