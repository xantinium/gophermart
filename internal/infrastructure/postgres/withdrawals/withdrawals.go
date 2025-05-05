package withdrawals

import (
	"context"
	"time"

	"github.com/xantinium/gophermart/internal/infrastructure/postgres/helpers"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres/withdrawals/gen"
	"github.com/xantinium/gophermart/internal/models"
)

func NewWithdrawalsTable(db gen.DBTX) *WithdrawalsTable {
	return &WithdrawalsTable{q: gen.New(db)}
}

type WithdrawalsTable struct {
	q *gen.Queries
}

func (t *WithdrawalsTable) CreateWithdrawal(ctx context.Context, userID int, order string, sum float32) error {
	now := time.Now()

	err := t.q.CreateWithdrawal(ctx, gen.CreateWithdrawalParams{
		OrderID: order,
		Sum:     sum,
		UserID:  int32(userID),
		Created: helpers.TimeToTimestamp(now),
		Updated: helpers.TimeToTimestamp(now),
	})

	return helpers.ConvertError(err)
}

func (t *WithdrawalsTable) GetWithdrawalsByUserID(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	withdrawals, err := t.q.GetWithdrawalsByUserID(ctx, int32(userID))
	if err != nil {
		return nil, helpers.ConvertError(err)
	}

	return convertWithdrawals(withdrawals), nil
}

func (t *WithdrawalsTable) GetTotalWithdrawnByUserID(ctx context.Context, userID int) (float32, error) {
	totalWithdrawn, err := t.q.GetTotalWithdrawnByUserID(ctx, int32(userID))
	if err != nil {
		return 0, helpers.ConvertError(err)
	}

	return totalWithdrawn, nil
}

func convertWithdrawals(withdrawals []gen.Withdrawal) []models.Withdrawal {
	result := make([]models.Withdrawal, len(withdrawals))

	for i := range withdrawals {
		result[i] = convertWithdrawal(withdrawals[i])
	}

	return result
}

func convertWithdrawal(withdrawal gen.Withdrawal) models.Withdrawal {
	return models.NewWithdrawal(
		int(withdrawal.ID),
		withdrawal.OrderID,
		withdrawal.Sum,
		int(withdrawal.UserID),
		helpers.TimestampToTime(withdrawal.Created),
		helpers.TimestampToTime(withdrawal.Updated),
	)
}
