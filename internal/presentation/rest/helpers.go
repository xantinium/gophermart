package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/tools"
)

type restHandler = func(ctx *gin.Context, server any) (int, any, error)

func register(group *gin.RouterGroup, method, path string, handler restHandler) {
	group.Handle(method, path, func(ctx *gin.Context) {
		status, response, err := handler(ctx, nil)
		respond(ctx, status, response, err)
	})
}

type restError struct {
	Error string `json:"error"`
}

func respond(ctx *gin.Context, status int, v any, err error) {
	var (
		response   []byte
		marshalErr error
	)

	if err == nil {
		response, marshalErr = tools.MarshalJSON(v)
	} else {
		response, marshalErr = tools.MarshalJSON(restError{Error: err.Error()})
	}
	if marshalErr != nil {
		status = http.StatusInternalServerError
		response = []byte(fmt.Sprintf(`{"error":"failed to marshal json: %v"}`, marshalErr))
	}

	ctx.Writer.WriteHeader(status)
	ctx.Writer.Write(response)
}

func runServer(server *http.Server) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	return errChan
}
