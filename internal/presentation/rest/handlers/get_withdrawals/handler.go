package getwithdrawalshandler

import (
	"errors"
	"log"
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

	withdrawals, err := server.GetUseCases().GetWithdrawals(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return http.StatusNoContent, response{}, nil
		default:
			logger.Errorf("failed to get withdrawals: %v", err)
			return http.StatusInternalServerError, response{}, err
		}
	}

	log.Println(withdrawals)

	res := make(response, len(withdrawals))
	for i := range withdrawals {
		res[i] = withdrawalItem{
			Order:       withdrawals[i].Order(),
			Sum:         withdrawals[i].Sum(),
			ProcessedAt: withdrawals[i].Created().Format(time.RFC3339),
		}
	}

	log.Println(res)

	return http.StatusOK, res, nil
}
