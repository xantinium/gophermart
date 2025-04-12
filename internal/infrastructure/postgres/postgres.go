package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/xantinium/gophermart/internal/models"
)

func init() {
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
}

// PostgresClientOptions опции для PostgreSQL клиента.
type PostgresClientOptions struct {
	PoolSize        int
	ConnMaxIdleTime time.Duration
	ConnMaxLifeTime time.Duration
}

// DefaultOptions опции по умолчанию для PostgreSQL клиента.
var DefaultOptions = PostgresClientOptions{
	PoolSize:        5,
	ConnMaxIdleTime: time.Second * 60,
	ConnMaxLifeTime: time.Minute * 5,
}

// NewPostgresClient создаёт новый клиент для работы с PostgreSQL.
func NewPostgresClient(ctx context.Context, connStr string, opts PostgresClientOptions) (*PostgresClient, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	client := &PostgresClient{
		db: db,
	}

	err = client.initTables(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return client, nil
}

// PostgresClient клиент для работы с PostgreSQL.
type PostgresClient struct {
	db *sql.DB
}

// Destroy уничтожает клиент, закрывая все соединения.
func (client *PostgresClient) Destroy() error {
	return client.db.Close()
}

func convertError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return models.ErrAlreadyExists
		}
	}

	return err
}
