package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/xantinium/gophermart/internal/tools"
	"github.com/xantinium/gophermart/internal/usecases"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func BindRequestBody(ctx *gin.Context, v any) error {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	defer ctx.Request.Body.Close()

	err = tools.UnmarshalJSON(reqBody, v)
	if err != nil {
		return err
	}

	return validate.StructCtx(ctx, v)
}

type RestServer interface {
	GetUseCases() *usecases.UseCases
}

type RestHandler[T any] interface {
	GetMethod() string
	Parse(ctx *gin.Context) (T, error)
	Handle(ctx *gin.Context, server RestServer, req T) (int, any, error)
}

type CustomRestHandler interface {
	GetMethod() string
	Handle(ctx *gin.Context, server RestServer)
}
