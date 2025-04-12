package registerhandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/models"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
)

func New() h {
	return h{}
}

type h struct{}

func (h) Handle(ctx *gin.Context, server handlers.RestServer, req request) (int, any, error) {
	err := server.GetUseCases().RegisterUser(ctx, req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAlreadyExists):
			return http.StatusConflict, response{}, fmt.Errorf("user with login %q already exists", req.Login)
		default:
			return http.StatusInternalServerError, response{}, err
		}
	}

	return http.StatusOK, response{}, nil
}
