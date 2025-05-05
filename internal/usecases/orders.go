package usecases

import (
	"context"
	"errors"

	"github.com/xantinium/gophermart/internal/models"
)

// CreateOrder создаёт заказ, если он
// ещё не был создан. Дополнительно возвращает
// признак создания нового заказа.
func (cases *UseCases) CreateOrder(ctx context.Context, number string, userID int) (bool, error) {
	err := cases.ordersRepo.CreateOrder(ctx, userID, number, models.OrderStatusNew, 0)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAlreadyExists):
			var order models.Order

			order, err = cases.ordersRepo.GetOrderByNumber(ctx, number)
			if err != nil {
				return false, err
			}

			if order.UserID() != userID {
				return false, models.ErrAlreadyExists
			}

			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (cases *UseCases) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return cases.ordersRepo.GetOrdersByUserID(ctx, userID)
}

func (cases *UseCases) GetAllOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	return cases.ordersRepo.GetOrders(ctx, limit, offset)
}

func (cases *UseCases) MarkOrderAsProcessing(ctx context.Context, number string) error {
	return cases.ordersRepo.UpdateOrder(ctx, number, models.OrderStatusProcessing, 0)
}

func (cases *UseCases) MarkOrderAsInvalid(ctx context.Context, number string) error {
	return cases.ordersRepo.UpdateOrder(ctx, number, models.OrderStatusInvalid, 0)
}

func (cases *UseCases) MarkOrderAsProcessed(ctx context.Context, number string, accrual float32) error {
	return cases.ordersRepo.UpdateOrder(ctx, number, models.OrderStatusProcessed, accrual)
}
