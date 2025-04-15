package usecases

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

// CreateOrder создаёт заказ, если он
// ещё не был создан. Дополнительно возвращает
// признак создания нового заказа.
func (cases *UseCases) CreateOrder(ctx context.Context, number string, userID int) (bool, error) {
	created, err := cases.ordersRepo.CreateOrder(ctx, userID, number, models.OrderStatusNew, 0)
	if err != nil {
		return false, err
	}

	return created, nil
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

func (cases *UseCases) MarkOrderAsProcessed(ctx context.Context, number string, accrual int) error {
	return cases.ordersRepo.UpdateOrder(ctx, number, models.OrderStatusProcessed, accrual)
}
