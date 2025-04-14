package ordersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type OrdersStorage interface {
	InsertOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual *int) (bool, error)
	FindOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error)
}
