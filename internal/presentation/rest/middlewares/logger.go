package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		duration := time.Since(start)
		status := ctx.Writer.Status()

		logger.Info(
			"api request",
			logger.Field{Name: "method", Value: method},
			logger.Field{Name: "status", Value: status},
			logger.Field{Name: "path", Value: path},
			logger.Field{Name: "query", Value: query},
			logger.Field{Name: "duration", Value: duration},
		)
	}
}
