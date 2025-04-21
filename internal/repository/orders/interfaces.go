package ordersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type OrdersStorage interface {
	InsertOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual float64) (bool, error)
	FindOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error)
	FindOrders(ctx context.Context, limit, offset int) ([]models.Order, error)
	UpdateOrder(ctx context.Context, number string, status models.OrderStatus, accrual float64) error
	SumAccrual(ctx context.Context, userID int) (float64, error)
}
