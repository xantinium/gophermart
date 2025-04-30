package ordersrepo

import (
	"context"

	"github.com/xantinium/gophermart/internal/models"
)

type OrdersStorage interface {
	CreateOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual float32) error
	GetOrderByNumber(ctx context.Context, number string) (models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error)
	GetOrdersByLimitAndOffset(ctx context.Context, limit, offset int) ([]models.Order, error)
	UpdateOrder(ctx context.Context, number string, status models.OrderStatus, accrual float32) error
	GetTotalAccrualByUserID(ctx context.Context, userID int) (float32, error)
}
