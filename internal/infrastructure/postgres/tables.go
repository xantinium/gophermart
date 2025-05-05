package postgres

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
)

func (client *PostgresClient) initTables(ctx context.Context) error {
	initializers := []func(context.Context) error{
		client.initUsersTable,
		client.initOrdersTable,
		client.initWithdrawalsTable,
	}

	for i := range initializers {
		err := initializers[i](ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *PostgresClient) initUsersTable(ctx context.Context) error {
	b := sqlbuilder.NewCreateTableBuilder()

	b.IfNotExists()
	b.CreateTable("users")
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("login", "VARCHAR(50)", "NOT NULL", "UNIQUE")
	b.Define("password_hash", "VARCHAR(64)", "NOT NULL")
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	query, args := b.Build()

	_, err := client.pool.Exec(ctx, query, args...)

	return err
}

func (client *PostgresClient) initOrdersTable(ctx context.Context) error {
	b := sqlbuilder.NewCreateTableBuilder()

	b.IfNotExists()
	b.CreateTable("orders")
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("number", "TEXT", "NOT NULL", "UNIQUE")
	b.Define("user_id", "INT", "REFERENCES users(id)")
	b.Define("status", "SMALLINT", "NOT NULL")
	b.Define("accrual", "REAL", "NOT NULL")
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	query, args := b.Build()

	_, err := client.pool.Exec(ctx, query, args...)

	return err
}

func (client *PostgresClient) initWithdrawalsTable(ctx context.Context) error {
	b := sqlbuilder.NewCreateTableBuilder()

	b.IfNotExists()
	b.CreateTable("withdrawals")
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("order_id", "TEXT", "NOT NULL", "UNIQUE")
	b.Define("sum", "REAL", "NOT NULL")
	b.Define("user_id", "INT", "REFERENCES users(id)")
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	query, args := b.Build()

	_, err := client.pool.Exec(ctx, query, args...)

	return err
}
