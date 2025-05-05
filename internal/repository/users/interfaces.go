package usersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type UsersStorage interface {
	CreateUser(ctx context.Context, login, passwordHash string) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}
