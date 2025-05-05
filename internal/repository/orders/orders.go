package ordersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

func New(storage OrdersStorage) *OrdersRepository {
	return &OrdersRepository{
		storage: storage,
	}
}

type OrdersRepository struct {
	storage OrdersStorage
}

func (repo *OrdersRepository) CreateOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual float32) error {
	return repo.storage.CreateOrder(ctx, userID, number, status, accrual)
}

func (repo *OrdersRepository) GetOrderByNumber(ctx context.Context, number string) (models.Order, error) {
	return repo.storage.GetOrderByNumber(ctx, number)
}

func (repo *OrdersRepository) GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	return repo.storage.GetOrdersByUserID(ctx, userID)
}

func (repo *OrdersRepository) GetOrders(ctx context.Context, limit, offset int) ([]models.Order, error) {
	return repo.storage.GetOrdersByLimitAndOffset(ctx, limit, offset)
}

func (repo *OrdersRepository) UpdateOrder(ctx context.Context, number string, status models.OrderStatus, accrual float32) error {
	return repo.storage.UpdateOrder(ctx, number, status, accrual)
}

func (repo *OrdersRepository) GetTotalAccrual(ctx context.Context, userID int) (float32, error) {
	return repo.storage.GetTotalAccrualByUserID(ctx, userID)
}
