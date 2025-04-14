package ordersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type OrdersStorage interface {
	FindOrderByNumber(ctx context.Context, number string) (models.Order, error)
	InsertOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual *int) error
	FindOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error)
}
