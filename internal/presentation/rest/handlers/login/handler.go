package loginhandler

import (
	"errors"
	"fmt"
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
	token, err := server.GetUseCases().AuthorizeUser(ctx, req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return http.StatusUnauthorized, response{}, fmt.Errorf("user does not exists")
		case errors.Is(err, models.ErrFailedToMatch):
			return http.StatusUnauthorized, response{}, fmt.Errorf("password does not match")
		default:
			logger.Errorf("failed to login: %v", err)
			return http.StatusInternalServerError, response{}, err
		}
	}

	tools.SetTokenCookie(ctx, token)

	return http.StatusOK, response{}, nil
}
