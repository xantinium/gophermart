package getordershandler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/models"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
	"github.com/xantinium/gophermart/internal/tools"
)

func New() h {
	return h{}
}

type h struct{}

func (h) Handle(ctx *gin.Context, server handlers.RestServer, req request) (int, any, error) {
	userID := tools.GetUserID(ctx)

	orders, err := server.GetUseCases().GetOrders(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return http.StatusNoContent, response{}, nil
		default:
			logger.Errorf("failed to get orders: %v", err)
			return http.StatusInternalServerError, response{}, err
		}
	}

	res := make(response, len(orders))
	for i := range orders {
		res[i] = orderItem{
			Number:     orders[i].Number(),
			Status:     orders[i].Status().String(),
			Accrual:    orders[i].Accrual(),
			UploadedAt: orders[i].Created().Format(time.RFC3339),
		}
	}

	return http.StatusOK, res, nil
}
