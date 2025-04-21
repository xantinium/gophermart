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
		poolSize:         opts.PoolSize,
		accrualHost:      opts.AccrualHost,
		useCases:         opts.UseCases,
		reader:           newOrdersReader(opts.UseCases, opts.PoolSize),
		requestLimiter:   tools.NewSemaphore(10),
		processingOrders: &ordersSet{set: map[string]struct{}{}},
	}
}

type WorkerPool struct {
	wg               sync.WaitGroup
	poolSize         int
	accrualHost      string
	useCases         *usecases.UseCases
	reader           *ordersReader
	requestLimiter   *tools.Semaphore
	processingOrders *ordersSet
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
					time.Sleep(time.Millisecond * 100)
					continue
				}

				if len(orders) == 0 {
					time.Sleep(time.Millisecond * 100)
					continue
				}

				for _, order := range orders {
					if w.processingOrders.Has(order.Number()) {
						time.Sleep(time.Millisecond * 100)
						continue
					}

					w.processingOrders.Add(order.Number())
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

					var (
						err  error
						info orderInfo
					)

					tools.DefaulRetrier.Exec(func() bool {
						info, err = w.requestOrdedInfo(ctx, order.Number())
						return err != nil
					})

					// TODO: возможно тут стоит устанавливать статус Invalid.
					if err != nil {
						logger.Errorf("failed to request order %q info: %v", order.Number(), err)
						w.processingOrders.Remove(order.Number())
						continue
					}

					if info.Status == order.Status() {
						logger.Infof("order has no changes")
						w.processingOrders.Remove(order.Number())
						continue
					}

					switch info.Status {
					case models.OrderStatusProcessing:
						err = w.useCases.MarkOrderAsProcessing(ctx, order.Number())
					case models.OrderStatusInvalid:
						err = w.useCases.MarkOrderAsInvalid(ctx, order.Number())
					case models.OrderStatusProcessed:
						err = w.useCases.MarkOrderAsProcessed(ctx, order.Number(), info.Accrual)
					default:
						logger.Infof("order %s still in processing", order.Number())
					}

					w.processingOrders.Remove(order.Number())

					if err != nil {
						logger.Infof("failed to process order %q", order.Number())
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
	Accrual float64
}

func (w *WorkerPool) requestOrdedInfo(ctx context.Context, number string) (orderInfo, error) {
	w.requestLimiter.Acquire()
	defer w.requestLimiter.Release()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/orders/%s", w.accrualHost, number), nil)
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
		Order   string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float64 `json:"accrual"`
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

type ordersSet struct {
	mx  sync.RWMutex
	set map[string]struct{}
}

func (orders *ordersSet) Has(number string) bool {
	orders.mx.RLock()
	defer orders.mx.RUnlock()

	_, exists := orders.set[number]
	return exists
}

func (orders *ordersSet) Add(number string) {
	orders.mx.Lock()
	defer orders.mx.Unlock()

	orders.set[number] = struct{}{}
}

func (orders *ordersSet) Remove(number string) {
	orders.mx.Lock()
	defer orders.mx.Unlock()

	delete(orders.set, number)
}
