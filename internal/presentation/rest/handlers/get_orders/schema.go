package getordershandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h) GetMethod() string {
	return http.MethodGet
}

type request struct{}

func (h) Parse(ctx *gin.Context) (request, error) {
	return request{}, nil
}

type orderItem struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

type response = []orderItem
