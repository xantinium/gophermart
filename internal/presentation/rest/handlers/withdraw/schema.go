package withdrawhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
)

func (h) GetMethod() string {
	return http.MethodPost
}

type request struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required,gt=0"`
}

func (h) Parse(ctx *gin.Context) (request, error) {
	var req request
	err := handlers.BindRequestBody(ctx, &req)
	return req, err
}

type response struct{}
