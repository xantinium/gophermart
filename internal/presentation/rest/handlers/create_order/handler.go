package createorderhandler

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/consts"
	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/models"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
	"github.com/xantinium/gophermart/internal/tools"
)

var errInvalidContentType = errors.New("invalid Content-Type header value")

func New() h {
	return h{}
}

type h struct{}

func (h) GetMethod() string {
	return http.MethodPost
}

func (h) Handle(ctx *gin.Context, server handlers.RestServer) {
	userID := tools.GetUserID(ctx)

	orderNum, err := parseRequest(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errInvalidContentType):
			ctx.AbortWithStatus(http.StatusBadRequest)
		case errors.Is(err, models.ErrInvalidOrderNum):
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		default:
			logger.Errorf("failed to create order: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	var created bool
	created, err = server.GetUseCases().CreateOrder(ctx, orderNum, userID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAlreadyExists):
			ctx.AbortWithStatus(http.StatusConflict)
		default:
			logger.Errorf("failed to create order: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	if !created {
		ctx.AbortWithStatus(http.StatusOK)
		return
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}

func parseRequest(ctx *gin.Context) (string, error) {
	if ctx.GetHeader(consts.HeaderContentType) != "text/plain" {
		return "", errInvalidContentType
	}

	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return "", err
	}

	orderNumber := strings.TrimSpace(string(reqBody))

	if ok := tools.CheckLuhn(orderNumber); !ok {
		return "", models.ErrInvalidOrderNum
	}

	return orderNumber, nil
}
