package withdrawalsrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type WithdrawalsStorage interface {
	FindWithdrawalsByUserID(ctx context.Context, userID int) ([]models.Withdrawal, error)
	InsertWithdrawal(ctx context.Context, userID int, order string, sum float64) error
	SumWithdrawn(ctx context.Context, userID int) (float64, error)
}
