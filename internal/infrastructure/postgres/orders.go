package postgres

import (
	"context"
	"time"

	"github.com/huandu/go-sqlbuilder"

	"github.com/xantinium/gophermart/internal/models"
)

func (client *PostgresClient) FindOrderByNumber(ctx context.Context, number string) (models.Order, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("*")
	b.From(OrdersTable)
	b.Where(b.Equal("number", number))
	b.Limit(1)

	query, args := b.Build()

	row := client.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return models.Order{}, convertError(row.Err())
	}

	var (
		orderID, orderUserID       int
		orderNumber                string
		orderStatus                models.OrderStatus
		orderCreated, orderUpdated time.Time
	)

	err := row.Scan(&orderID, &orderNumber, &orderUserID, &orderStatus, &orderCreated, &orderUpdated)
	if err != nil {
		return models.Order{}, convertError(err)
	}

	return models.NewOrder(orderID, orderNumber, orderUserID, orderStatus, orderCreated, orderUpdated), nil
}

func (client *PostgresClient) InsertOrder(ctx context.Context, userID int, number string, status models.OrderStatus) error {
	now := time.Now()
	b := sqlbuilder.NewInsertBuilder()

	b.InsertInto(OrdersTable)
	b.Cols("number", "user_id", "status", "created", "updated")
	b.Values(number, userID, status, now, now)

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}
