package postgres

import (
	"context"
	"fmt"

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
	b.CreateTable(UsersTable)
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("login", "VARCHAR(50)", "NOT NULL", "UNIQUE")
	b.Define("password_hash", "VARCHAR(64)", "NOT NULL")
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}

func (client *PostgresClient) initOrdersTable(ctx context.Context) error {
	b := sqlbuilder.NewCreateTableBuilder()

	b.IfNotExists()
	b.CreateTable(OrdersTable)
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("number", "TEXT", "NOT NULL", "UNIQUE")
	b.Define("user_id", "INT", fmt.Sprintf("REFERENCES %s(id)", UsersTable))
	b.Define("status", "SMALLINT", "NOT NULL")
	b.Define("accrual", "REAL", "NOT NULL")
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}

func (client *PostgresClient) initWithdrawalsTable(ctx context.Context) error {
	b := sqlbuilder.NewCreateTableBuilder()

	b.IfNotExists()
	b.CreateTable(WithdrawalsTable)
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("order", "TEXT", "NOT NULL", "UNIQUE")
	b.Define("sum", "INT", "NOT NULL")
	b.Define("user_id", "INT", fmt.Sprintf("REFERENCES %s(id)", UsersTable))
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}
