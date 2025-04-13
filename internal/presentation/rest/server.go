package rest

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/gophermart/internal/logger"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers/loginhandler"
	"github.com/xantinium/gophermart/internal/presentation/rest/handlers/registerhandler"
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

	apiGroup := engine.Group("/api")
	apiGroup.Use(middlewares.CompressMiddleware(), middlewares.LoggerMiddleware())

	registerPublicHandlers(server, apiGroup.Group(""))
	registerPrivateHandlers(server, apiGroup.Group(""))

	return server
}

func registerPublicHandlers(server *Server, rootGroup *gin.RouterGroup) {
	userGroup := rootGroup.Group("/user")
	{
		register(server, userGroup, "/register", registerhandler.New())
		register(server, userGroup, "/login", loginhandler.New())
	}
}

func registerPrivateHandlers(server *Server, rootGroup *gin.RouterGroup) {
	rootGroup.Use(middlewares.AuthMiddleware(server.useCases))
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
