package rest

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
	handlers "github.com/xantinium/gophermart/internal/presentation/rest/handlers/login"
	"github.com/xantinium/gophermart/internal/presentation/rest/middlewares"
)

type ServerOptions struct {
	IsDev bool
	Addr  string
}

func NewServer(opts ServerOptions) *Server {
	if !opts.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middlewares.RecoveryMiddleware(), middlewares.LoggerMiddleware())

	api := engine.Group("/api")

	privateRoutes := api.Group("")
	privateRoutes.Use()

	userGroup := privateRoutes.Group("/user")
	{
		register(userGroup, http.MethodPost, "/login", handlers.LoginHandler)
	}

	return &Server{
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: engine,
		},
	}
}

type Server struct {
	wg     sync.WaitGroup
	server *http.Server
}

func (s *Server) Run(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		select {
		case <-ctx.Done():
			return
		case err := <-runServer(s.server):
			if err != nil {
				panic(err)
			}
		}
	}()
}

func (s *Server) Wait() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Errorf("failed to gracefully shutdown server: %w", err)
	}
}
