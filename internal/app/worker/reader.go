package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/xantinium/gophermart/internal/models"
	"github.com/xantinium/gophermart/internal/usecases"
)

func newOrdersReader(useCases *usecases.UseCases, batchSize int) *ordersReader {
	return &ordersReader{
		offset:    0,
		batchSize: batchSize,
		useCases:  useCases,
	}
}

type ordersReader struct {
	mx        sync.Mutex
	offset    int
	batchSize int
	useCases  *usecases.UseCases
}

func (reader *ordersReader) ReadOrders(ctx context.Context) ([]models.Order, error) {
	reader.mx.Lock()
	defer reader.mx.Unlock()

	orders, err := reader.useCases.GetAllOrders(ctx, reader.batchSize, reader.offset)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			reader.offset = 0
			return []models.Order{}, nil
		default:
			return nil, err
		}
	}

	if len(orders) < reader.batchSize {
		reader.offset = 0
	} else {
		reader.offset += reader.batchSize
	}

	return orders, err
}
