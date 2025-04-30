package withdrawalsrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type WithdrawalsStorage interface {
	GetWithdrawalsByUserID(ctx context.Context, userID int) ([]models.Withdrawal, error)
	CreateWithdrawal(ctx context.Context, userID int, order string, sum float32) error
	GetTotalWithdrawnByUserID(ctx context.Context, userID int) (float32, error)
}
