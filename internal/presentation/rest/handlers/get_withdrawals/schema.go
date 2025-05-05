package getwithdrawalshandler

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

type withdrawalItem struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type response = []withdrawalItem
