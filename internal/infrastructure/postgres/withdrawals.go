package postgres

import (
	"context"
	"time"

	"github.com/huandu/go-sqlbuilder"

	"github.com/xantinium/gophermart/internal/models"
)

func (client *PostgresClient) InsertWithdrawal(ctx context.Context, userID int, order string, sum float64) error {
	now := time.Now()
	b := sqlbuilder.NewInsertBuilder()

	b.InsertInto(WithdrawalsTable)
	b.Cols("order_id", "sum", "user_id", "created", "updated")
	b.Values(order, sum, userID, now, now)

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}

func (client *PostgresClient) FindWithdrawalsByUserID(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("*")
	b.From(WithdrawalsTable)
	b.Where(b.Equal("user_id", userID))

	query, args := b.Build()

	rows, err := client.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, convertError(err)
	}
	defer rows.Close()

	withdrawals := make([]models.Withdrawal, 0)

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, convertError(err)
		}

		var (
			withdrawalID, withdrawalUserID       int
			withdrawalSum                        float64
			withdrawalOrder                      string
			withdrawalCreated, withdrawalUpdated time.Time
		)

		err = rows.Scan(&withdrawalID, &withdrawalOrder, &withdrawalSum, &withdrawalUserID, &withdrawalCreated, &withdrawalUpdated)
		if err != nil {
			return nil, convertError(err)
		}

		withdrawals = append(withdrawals, models.NewWithdrawal(withdrawalID, withdrawalOrder, withdrawalSum, withdrawalUserID, withdrawalCreated, withdrawalUpdated))
	}

	if len(withdrawals) == 0 {
		return nil, models.ErrNotFound
	}

	return withdrawals, nil
}

func (client *PostgresClient) SumWithdrawn(ctx context.Context, userID int) (float64, error) {
	b := sqlbuilder.NewSelectBuilder()

	b.Select("COALESCE(SUM(sum), 0) as total_sum")
	b.From(WithdrawalsTable)
	b.Where(b.Equal("user_id", userID))

	query, args := b.Build()

	row := client.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return 0, convertError(row.Err())
	}

	var totalSum float64
	err := row.Scan(&totalSum)

	return totalSum, convertError(err)
}
