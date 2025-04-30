package usecases

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

func (cases *UseCases) CreateWithdrawal(ctx context.Context, order string, sum float32, userID int) error {
	balance, err := cases.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.AvaliableAccrual < sum {
		return models.ErrInsufficientBalance
	}

	return cases.withdrawalsRepo.CreateWithdrawal(ctx, userID, order, sum)
}

func (cases *UseCases) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return cases.withdrawalsRepo.GetWithdrawalsByUserID(ctx, userID)
}
