package usecases

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

// CreateOrder создаёт заказ, если он
// ещё не был создан. Дополнительно возвращает
// признак создания нового заказа.
func (cases *UseCases) CreateOrder(ctx context.Context, number string, userID int) (bool, error) {
	created, err := cases.ordersRepo.CreateOrder(ctx, userID, number, models.OrderStatusNew, nil)
	if err != nil {
		return false, err
	}

	return created, nil
}

func (cases *UseCases) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return cases.ordersRepo.GetOrders(ctx, userID)
}
