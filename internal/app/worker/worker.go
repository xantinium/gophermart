package worker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/models"
	"github.com/xantinium/gophermart/internal/tools"
	"github.com/xantinium/gophermart/internal/usecases"
)

type WorkerPoolOptions struct {
	PoolSize    int
	AccrualHost string
	UseCases    *usecases.UseCases
}

func NewWorkerPool(opts WorkerPoolOptions) *WorkerPool {
	return &WorkerPool{
		poolSize:       opts.PoolSize,
		accrualHost:    opts.AccrualHost,
		useCases:       opts.UseCases,
		reader:         newOrdersReader(opts.UseCases, opts.PoolSize),
		requestLimiter: tools.NewSemaphore(10),
	}
}

type WorkerPool struct {
	wg             sync.WaitGroup
	poolSize       int
	accrualHost    string
	useCases       *usecases.UseCases
	reader         *ordersReader
	requestLimiter *tools.Semaphore
}

func (w *WorkerPool) Run(ctx context.Context) {
	ordersChan := make(chan models.Order, w.poolSize)

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				orders, err := w.reader.ReadOrders(ctx)
				if err != nil {
					logger.Errorf("failed to read orders: %v", err)
					continue
				}

				if len(orders) == 0 {
					time.Sleep(time.Millisecond * 100)
					continue
				}

				for _, order := range orders {
					ordersChan <- order
				}
			}
		}
	}()

	for range w.poolSize {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case order, ok := <-ordersChan:
					if !ok {
						return
					}

					info, err := w.requestOrdedInfo(ctx, order.Number())
					if err != nil {
						logger.Errorf("failed to request order info: %v", err)
						continue
					}

					if info.Status == order.Status() {
						logger.Infof("order has no changes")
						continue
					}

					switch info.Status {
					case models.OrderStatusProcessing:
						w.useCases.MarkOrderAsProcessing(ctx, order.Number())
					case models.OrderStatusInvalid:
						w.useCases.MarkOrderAsInvalid(ctx, order.Number())
					case models.OrderStatusProcessed:
						w.useCases.MarkOrderAsProcessed(ctx, order.Number(), info.Accrual)
					default:
						logger.Infof("order still in processing")
					}
				}
			}
		}()
	}
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}

type orderInfo struct {
	Order   string
	Status  models.OrderStatus
	Accrual int
}

func (w *WorkerPool) requestOrdedInfo(ctx context.Context, number string) (orderInfo, error) {
	w.requestLimiter.Acquire()
	defer w.requestLimiter.Release()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/api/orders/%s", w.accrualHost, number), nil)
	if err != nil {
		return orderInfo{}, err
	}

	rawResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return orderInfo{}, err
	}
	defer rawResp.Body.Close()

	var respBytes []byte
	respBytes, err = io.ReadAll(rawResp.Body)
	if err != nil {
		return orderInfo{}, err
	}

	var resp struct {
		Order   string `json:"order"`
		Status  string `json:"status"`
		Accrual int    `json:"accrual"`
	}
	err = tools.UnmarshalJSON(respBytes, &resp)
	if err != nil {
		return orderInfo{}, err
	}

	var status models.OrderStatus
	status, err = parseOrderStatus(resp.Status)
	if err != nil {
		return orderInfo{}, err
	}

	return orderInfo{
		Order:   resp.Order,
		Status:  status,
		Accrual: resp.Accrual,
	}, nil
}

func parseOrderStatus(maybeStatus string) (models.OrderStatus, error) {
	switch maybeStatus {
	case "NEW":
		return models.OrderStatusNew, nil
	case "PROCESSING":
		return models.OrderStatusProcessing, nil
	case "INVALID":
		return models.OrderStatusInvalid, nil
	case "PROCESSED":
		return models.OrderStatusProcessed, nil
	default:
		return 0, fmt.Errorf("unknown status %q", maybeStatus)
	}
}
