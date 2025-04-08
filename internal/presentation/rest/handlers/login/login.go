package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context, server any) (int, any, error) {
	req, err := parseReq(ctx)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, req, nil
}
