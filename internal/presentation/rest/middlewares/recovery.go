package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(recoveryWriter{})
}

type recoveryWriter struct{}

func (w recoveryWriter) Write(p []byte) (int, error) {
	logger.Error("panic", logger.Field{Name: "err", Value: string(p)})
	return len(p), nil
}
