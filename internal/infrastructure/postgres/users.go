package postgres

import (
	"context"
	"time"

	"github.com/huandu/go-sqlbuilder"

	"github.com/xantinium/gophermart/internal/models"
)

func (client *PostgresClient) InsertUser(ctx context.Context, login, passwordHash string) error {
	now := time.Now()
	b := sqlbuilder.NewInsertBuilder()

	b.InsertInto(UsersTable)
	b.Cols("login", "password_hash", "created", "updated")
	b.Values(login, passwordHash, now, now)

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}

func (client *PostgresClient) FindUserByLogin(ctx context.Context, login string) (models.User, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("*")
	b.From(UsersTable)
	b.Where(b.Equal("login", login))
	b.Limit(1)

	query, args := b.Build()

	row := client.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return models.User{}, convertError(row.Err())
	}

	var (
		userID                      int
		userLogin, userPasswordHash string
		userCreated, userUpdated    time.Time
	)

	err := row.Scan(&userID, &userLogin, &userPasswordHash, &userCreated, &userUpdated)
	if err != nil {
		return models.User{}, convertError(err)
	}

	return models.NewUser(userID, userLogin, userPasswordHash, userCreated, userUpdated), nil
}
