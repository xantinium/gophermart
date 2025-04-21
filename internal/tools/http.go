package tools

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/consts"
)

const (
	tokenCookieName = "token"
	userIDKey       = "user_id"
)

var ErrNoUserID = errors.New("no user_id in context")

func GetTokenCookie(ctx *gin.Context) (string, error) {
	return ctx.Cookie(tokenCookieName)
}

func SetTokenCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(tokenCookieName, token, 0, "", "", false, true)
}

func GetUserID(ctx *gin.Context) int {
	userID := ctx.GetInt(userIDKey)
	if userID == 0 {
		panic(ErrNoUserID)
	}

	return userID
}

func SetUserID(ctx *gin.Context, userID int) {
	ctx.Set(userIDKey, userID)
}

func WriteJSON(ctx *gin.Context, statusCode int, json []byte) {
	ctx.Header(consts.HeaderContentType, "application/json; charset=utf-8")
	ctx.String(statusCode, string(json))
}
