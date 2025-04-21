package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/huandu/go-sqlbuilder"

	"github.com/xantinium/gophermart/internal/models"
)

// InsertOrder добавляет заказ в таблицу заказов.
// Дополнительно возвращает признак создания заказа.
//
// TODO: скорее всего может быть решено одним запросом.
func (client *PostgresClient) InsertOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual float64) (bool, error) {
	now := time.Now()
	b := sqlbuilder.NewInsertBuilder()

	b.InsertInto(OrdersTable)
	b.Cols("number", "user_id", "status", "accrual", "created", "updated")
	b.Values(number, userID, status, accrual, now, now)

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)
	if err != nil {
		err = convertError(err)

		switch {
		case errors.Is(err, models.ErrAlreadyExists):
			var orderUserID int
			orderUserID, err = client.findUserIDByNumber(ctx, number)
			if err != nil {
				return false, err
			}

			if orderUserID != userID {
				return false, models.ErrAlreadyExists
			}

			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (client *PostgresClient) FindOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("*")
	b.From(OrdersTable)
	b.Where(b.Equal("user_id", userID))

	query, args := b.Build()

	rows, err := client.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, convertError(err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, convertError(err)
		}

		var (
			orderID, orderUserID       int
			orderAccrual               float64
			orderNumber                string
			orderStatus                models.OrderStatus
			orderCreated, orderUpdated time.Time
		)

		err = rows.Scan(&orderID, &orderNumber, &orderUserID, &orderStatus, &orderAccrual, &orderCreated, &orderUpdated)
		if err != nil {
			return nil, convertError(err)
		}

		orders = append(orders, models.NewOrder(orderID, orderNumber, orderUserID, orderStatus, orderAccrual, orderCreated, orderUpdated))
	}

	if len(orders) == 0 {
		return nil, models.ErrNotFound
	}

	return orders, nil
}

func (client *PostgresClient) findUserIDByNumber(ctx context.Context, number string) (int, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("user_id")
	b.From(OrdersTable)
	b.Where(b.Equal("number", number))
	b.OrderBy("created")
	b.Desc()

	query, args := b.Build()

	row := client.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return 0, convertError(row.Err())
	}

	var userID int
	err := row.Scan(&userID)

	return userID, convertError(err)
}

func (client *PostgresClient) FindOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("*")
	b.From(OrdersTable)
	b.Where(b.In("status", models.OrderStatusNew, models.OrderStatusProcessing))
	b.OrderBy("id")
	b.Asc()
	b.Limit(limit)
	b.Offset(offset)

	query, args := b.Build()

	rows, err := client.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, convertError(err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, convertError(err)
		}

		var (
			orderID, orderUserID       int
			orderAccrual               float64
			orderNumber                string
			orderStatus                models.OrderStatus
			orderCreated, orderUpdated time.Time
		)

		err = rows.Scan(&orderID, &orderNumber, &orderUserID, &orderStatus, &orderAccrual, &orderCreated, &orderUpdated)
		if err != nil {
			return nil, convertError(err)
		}

		orders = append(orders, models.NewOrder(orderID, orderNumber, orderUserID, orderStatus, orderAccrual, orderCreated, orderUpdated))
	}

	if len(orders) == 0 {
		return nil, models.ErrNotFound
	}

	return orders, nil
}

func (client *PostgresClient) UpdateOrder(ctx context.Context, number string, status models.OrderStatus, accrual float64) error {
	now := time.Now()
	b := sqlbuilder.NewUpdateBuilder()

	b.Update(OrdersTable)
	b.Where(b.Equal("number", number))
	b.Set(b.Assign("status", status), b.Assign("accrual", accrual), b.Assign("updated", now))

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (client *PostgresClient) SumAccrual(ctx context.Context, userID int) (float64, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("COALESCE(SUM(accrual), 0) as total_accrual")
	b.From(OrdersTable)
	b.Where(b.Equal("user_id", userID), b.Equal("status", models.OrderStatusProcessed))

	query, args := b.Build()

	row := client.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return 0, convertError(row.Err())
	}

	var totalAccrual float64
	err := row.Scan(&totalAccrual)

	return totalAccrual, convertError(err)
}
