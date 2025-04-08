package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
)

type loginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func parseReq(ctx *gin.Context) (loginRequest, error) {
	var req loginRequest
	err := handlers.BindRequestBody(ctx, &req)
	return req, err
}
