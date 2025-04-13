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

func (repo *OrdersRepository) GetOrderByNumber(ctx context.Context, number string) (models.Order, error) {
	return repo.storage.FindOrderByNumber(ctx, number)
}

func (repo *OrdersRepository) CreateOrder(ctx context.Context, userID int, number string, status models.OrderStatus) error {
	return repo.storage.InsertOrder(ctx, userID, number, status)
}
