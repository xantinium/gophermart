package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/tools"
	"github.com/xantinium/gophermart/internal/usecases"
)

func AuthMiddleware(cases *usecases.UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := tools.GetTokenCookie(ctx)
		if err != nil {
			tools.WriteJSON(ctx, http.StatusUnauthorized, fmt.Appendf(nil, `{"error":"%v"}`, err))
			return
		}

		err = cases.VerifyUserAuthorization(ctx, token)
		if err != nil {
			tools.WriteJSON(ctx, http.StatusUnauthorized, fmt.Appendf(nil, `{"error":"%v"}`, err))
			return
		}

		ctx.Next()
	}
}
