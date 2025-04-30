package helpers

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xantinium/gophermart/internal/models"
)

func ConvertError(err error) error {
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
