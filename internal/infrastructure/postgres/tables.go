package postgres

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
)

func (client *PostgresClient) initTables(ctx context.Context) error {
	initializers := []func(context.Context) error{
		client.initUsersTable,
		client.initOrdersTable,
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
	b := sqlbuilder.PostgreSQL.NewCreateTableBuilder()

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
	b := sqlbuilder.PostgreSQL.NewCreateTableBuilder()

	b.IfNotExists()
	b.CreateTable(OrdersTable)
	b.Define("id", "SERIAL", "PRIMARY KEY")
	b.Define("number", "STRING", "NOT NULL", "UNIQUE")
	b.Define("user_id", "VARCHAR(50)", "NOT NULL", "UNIQUE")
	b.Define("status", "SMALLINT", "NOT NULL")
	b.Define("created", "TIMESTAMP", "NOT NULL")
	b.Define("updated", "TIMESTAMP", "NOT NULL")

	b.Define("CONSTRAINT", "fk_user")
	b.Define("FOREIGN KEY", "user_id")
	b.Define("REFERENCES", UsersTable+".id")
	b.Define("ON DELETE", "CASCADE")

	query, args := b.Build()

	_, err := client.db.ExecContext(ctx, query, args...)

	return convertError(err)
}
