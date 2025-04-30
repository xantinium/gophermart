package getbalancehandler

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

type response struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
