package getbalancehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
	"github.com/xantinium/gophermart/internal/tools"
)

func New() h {
	return h{}
}

type h struct{}

func (h) Handle(ctx *gin.Context, server handlers.RestServer, req request) (int, any, error) {
	userID := tools.GetUserID(ctx)

	balance, err := server.GetUseCases().GetUserBalance(ctx, userID)
	if err != nil {
		logger.Errorf("failed to get balance: %v", err)
		return http.StatusInternalServerError, response{}, err
	}

	res := response{
		Current:   balance.AvaliableAccrual,
		Withdrawn: balance.TotalWithdrawn,
	}

	return http.StatusOK, res, nil
}
