package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
	login_handler "github.com/xantinium/gophermart/internal/presentation/rest/handlers/login"
	register_handler "github.com/xantinium/gophermart/internal/presentation/rest/handlers/register"
	"github.com/xantinium/gophermart/internal/presentation/rest/middlewares"
	"github.com/xantinium/gophermart/internal/usecases"
)

type ServerOptions struct {
	IsDev    bool
	Addr     string
	UseCases *usecases.UseCases
}

func NewServer(opts ServerOptions) *Server {
	if !opts.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(middlewares.RecoveryMiddleware())

	server := &Server{
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: engine,
		},
		useCases:      opts.UseCases,
		tokensCleaner: NewTokensCleaner(opts.UseCases),
	}

	api := engine.Group("/api")
	api.Use(middlewares.LoggerMiddleware())

	registerPublicHandlers(server, api.Group(""))
	registerPrivateHandlers(server, api.Group(""))

	return server
}

func registerPublicHandlers(server *Server, root *gin.RouterGroup) {
	userGroup := root.Group("/user")
	{
		register(server, userGroup, "/register", register_handler.New())
		register(server, userGroup, "/login", login_handler.New())
	}
}

func registerPrivateHandlers(server *Server, root *gin.RouterGroup) {
	root.Use(middlewares.AuthMiddleware(server.useCases))
}

type Server struct {
	wg            sync.WaitGroup
	server        *http.Server
	useCases      *usecases.UseCases
	tokensCleaner *TokensCleaner
}

func (s *Server) Run(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		select {
		case <-ctx.Done():
			return
		case err := <-runServer(s.server):
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	s.tokensCleaner.Run(ctx)
}

func (s *Server) Wait() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Errorf("failed to gracefully shutdown server: %v", err)
	}

	s.tokensCleaner.Wait()
}

func (s *Server) GetUseCases() *usecases.UseCases {
	return s.useCases
}
