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
	existingOrder, err := cases.ordersRepo.GetOrderByNumber(ctx, number)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			err = cases.ordersRepo.CreateOrder(ctx, userID, number, models.OrderStatusNew, nil)
			if err != nil {
				return false, err
			}

			return true, nil
		default:
			return false, err
		}
	}

	if existingOrder.UserID() != userID {
		return false, models.ErrAlreadyExists
	}

	return false, nil
}

func (cases *UseCases) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return cases.ordersRepo.GetOrders(ctx, userID)
}
