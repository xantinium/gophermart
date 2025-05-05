package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/tools"
	"github.com/xantinium/gophermart/internal/usecases"
)

func AuthMiddleware(cases *usecases.UseCases) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				if r == tools.ErrNoUserID {
					ctx.AbortWithStatus(http.StatusUnauthorized)
					return
				}

				panic(r)
			}
		}()

		var (
			err    error
			token  string
			userID int
		)

		token, err = tools.GetTokenCookie(ctx)
		if err != nil {
			logger.Errorf("failed to authorize user: %v", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err = cases.VerifyUserAuthorization(ctx, token)
		if err != nil {
			logger.Errorf("failed to authorize user: %v", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tools.SetUserID(ctx, userID)

		ctx.Next()
	}
}
