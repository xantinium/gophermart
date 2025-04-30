package withdrawalsrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

func New(storage WithdrawalsStorage) *WithdrawalsRepository {
	return &WithdrawalsRepository{
		storage: storage,
	}
}

type WithdrawalsRepository struct {
	storage WithdrawalsStorage
}

func (repo *WithdrawalsRepository) CreateWithdrawal(ctx context.Context, userID int, order string, sum float32) error {
	return repo.storage.CreateWithdrawal(ctx, userID, order, sum)
}

func (repo *WithdrawalsRepository) GetWithdrawalsByUserID(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return repo.storage.GetWithdrawalsByUserID(ctx, userID)
}

func (repo *WithdrawalsRepository) GetTotalWithdrawn(ctx context.Context, userID int) (float32, error) {
	return repo.storage.GetTotalWithdrawnByUserID(ctx, userID)
}
