package tools

import (
	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/consts"
)

const tokenCookieName = "token"

func GetTokenCookie(ctx *gin.Context) (string, error) {
	return ctx.Cookie(tokenCookieName)
}

func SetTokenCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(tokenCookieName, token, 0, "", "", true, true)
}

func WriteJSON(ctx *gin.Context, statusCode int, json []byte) {
	ctx.Header(consts.HeaderContentType, "application/json; charset=utf-8")
	ctx.String(statusCode, string(json))
}
