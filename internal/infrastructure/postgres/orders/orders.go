package orders

import (
	"context"
	"time"

	"github.com/xantinium/gophermart/internal/infrastructure/postgres/helpers"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres/orders/gen"
	"github.com/xantinium/gophermart/internal/models"
)

func NewOrdersTable(db gen.DBTX) *OrdersTable {
	return &OrdersTable{q: gen.New(db)}
}

type OrdersTable struct {
	q *gen.Queries
}

func (t *OrdersTable) CreateOrder(ctx context.Context, userID int, number string, status models.OrderStatus, accrual float32) error {
	now := time.Now()

	err := t.q.CreateOrder(ctx, gen.CreateOrderParams{
		Number:  number,
		UserID:  int32(userID),
		Status:  int16(status),
		Accrual: accrual,
		Created: helpers.TimeToTimestamp(now),
		Updated: helpers.TimeToTimestamp(now),
	})

	return helpers.ConvertError(err)
}

func (t *OrdersTable) GetOrderByNumber(ctx context.Context, number string) (models.Order, error) {
	order, err := t.q.GetOrderByNumber(ctx, number)
	if err != nil {
		return models.Order{}, helpers.ConvertError(err)
	}

	return convertOrder(order), nil
}

func (t *OrdersTable) GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := t.q.GetOrdersByUserID(ctx, int32(userID))
	if err != nil {
		return nil, helpers.ConvertError(err)
	}

	return convertOrders(orders), nil
}

func (t *OrdersTable) GetOrdersByLimitAndOffset(ctx context.Context, limit, offset int) ([]models.Order, error) {
	orders, err := t.q.GetOrdersByLimitAndOffset(ctx, gen.GetOrdersByLimitAndOffsetParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Statuses: []int16{int16(models.OrderStatusNew), int16(models.OrderStatusProcessing)},
	})
	if err != nil {
		return nil, helpers.ConvertError(err)
	}

	return convertOrders(orders), nil
}

func (t *OrdersTable) UpdateOrder(ctx context.Context, number string, status models.OrderStatus, accrual float32) error {
	now := time.Now()

	err := t.q.UpdateOrder(ctx, gen.UpdateOrderParams{
		Status:  int16(status),
		Accrual: float32(accrual),
		Updated: helpers.TimeToTimestamp(now),
		Number:  number,
	})

	return helpers.ConvertError(err)
}

func (t *OrdersTable) GetTotalAccrualByUserID(ctx context.Context, userID int) (float32, error) {
	totalAccrual, err := t.q.GetTotalAccrualByUserID(ctx, gen.GetTotalAccrualByUserIDParams{
		UserID: int32(userID),
		Status: int16(models.OrderStatusProcessed),
	})
	if err != nil {
		return 0, helpers.ConvertError(err)
	}

	return totalAccrual, nil
}

func convertOrders(orders []gen.Order) []models.Order {
	result := make([]models.Order, len(orders))

	for i := range orders {
		result[i] = convertOrder(orders[i])
	}

	return result
}

func convertOrder(order gen.Order) models.Order {
	return models.NewOrder(
		int(order.ID),
		order.Number,
		int(order.UserID),
		models.OrderStatus(order.Status),
		order.Accrual,
		helpers.TimestampToTime(order.Created),
		helpers.TimestampToTime(order.Updated),
	)
}
