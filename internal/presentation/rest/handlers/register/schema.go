package register_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
)

func (h) GetMethod() string {
	return http.MethodPost
}

type request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h) Parse(ctx *gin.Context) (request, error) {
	var req request
	err := handlers.BindRequestBody(ctx, &req)
	return req, err
}

type response struct{}
