package usersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type UsersStorage interface {
	InsertUser(ctx context.Context, login, passwordHash string) error
	FindUserByLogin(ctx context.Context, login string) (models.User, error)
}
