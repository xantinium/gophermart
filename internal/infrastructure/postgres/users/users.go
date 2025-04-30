package users

import (
	"context"
	"time"

	"github.com/xantinium/gophermart/internal/infrastructure/postgres/helpers"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres/users/gen"
	"github.com/xantinium/gophermart/internal/models"
)

func NewUsersTable(db gen.DBTX) *UsersTable {
	return &UsersTable{q: gen.New(db)}
}

type UsersTable struct {
	q *gen.Queries
}

func (t *UsersTable) CreateUser(ctx context.Context, login, passwordHash string) error {
	now := time.Now()

	err := t.q.CreateUser(ctx, gen.CreateUserParams{
		Login:        login,
		PasswordHash: passwordHash,
		Created:      helpers.TimeToTimestamp(now),
		Updated:      helpers.TimeToTimestamp(now),
	})

	return helpers.ConvertError(err)
}

func (t *UsersTable) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	user, err := t.q.GetUserByLogin(ctx, login)
	if err != nil {
		helpers.ConvertError(err)
	}

	return models.NewUser(
		int(user.ID),
		user.Login,
		user.PasswordHash,
		helpers.TimestampToTime(user.Created),
		helpers.TimestampToTime(user.Created),
	), nil
}
