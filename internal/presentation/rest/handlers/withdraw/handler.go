package withdrawhandler

import (
	"errors"
	"net/http"

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

	if ok := tools.CheckLuhn(req.Order); !ok {
		return http.StatusUnprocessableEntity, response{}, models.ErrInvalidOrderNum
	}

	err := server.GetUseCases().CreateWithdrawal(ctx, req.Order, req.Sum, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInsufficientBalance):
			return http.StatusPaymentRequired, response{}, err
		default:
			logger.Errorf("failed to withdraw: %v", err)
			return http.StatusInternalServerError, response{}, err
		}
	}

	return http.StatusOK, response{}, nil
}
