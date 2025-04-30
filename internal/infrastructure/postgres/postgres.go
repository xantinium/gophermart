package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/xantinium/gophermart/internal/infrastructure/postgres/orders"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres/users"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres/withdrawals"
)

// NewPostgresClient создаёт новый клиент для работы с PostgreSQL.
func NewPostgresClient(ctx context.Context, connStr string) (*PostgresClient, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	client := &PostgresClient{
		pool:             pool,
		UsersTable:       users.NewUsersTable(pool),
		OrdersTable:      orders.NewOrdersTable(pool),
		WithdrawalsTable: withdrawals.NewWithdrawalsTable(pool),
	}

	// sqlc не генерирует код для создания таблиц.
	// Скорее всего, предполагается, что это не входит
	// в зону ответственности приложения.
	err = client.initTables(ctx)
	if err != nil {
		pool.Close()

		return nil, err
	}

	return client, nil
}

// PostgresClient клиент для работы с PostgreSQL.
type PostgresClient struct {
	pool *pgxpool.Pool

	*users.UsersTable
	*orders.OrdersTable
	*withdrawals.WithdrawalsTable
}

// Destroy уничтожает клиент, закрывая все соединения.
func (client *PostgresClient) Destroy() {
	client.pool.Close()
}
