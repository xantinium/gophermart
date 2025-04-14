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

func (repo *OrdersRepository) CreateOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual *int) (bool, error) {
	return repo.storage.InsertOrder(ctx, userID, number, status, accrual)
}

func (repo *OrdersRepository) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return repo.storage.FindOrdersByUserID(ctx, userID)
}
