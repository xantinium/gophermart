package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/presentation/rest/handlers"
	"github.com/xantinium/gophermart/internal/tools"
)

func register[T any](server handlers.RestServer, group *gin.RouterGroup, path string, handler handlers.RestHandler[T]) {
	group.Handle(handler.GetMethod(), path, func(ctx *gin.Context) {
		var (
			err      error
			req      T
			status   int
			response any
		)

		req, err = handler.Parse(ctx)
		if err != nil {
			respond(ctx, http.StatusBadRequest, nil, err)
			return
		}

		status, response, err = handler.Handle(ctx, server, req)
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
		response = fmt.Appendf(nil, `{"error":"failed to marshal json: %v"}`, marshalErr)
	}

	tools.WriteJSON(ctx, status, response)
}

func runServer(server *http.Server) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	return errChan
}
