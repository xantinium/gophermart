package app

import (
	"context"

	"github.com/xantinium/gophermart/internal/presentation/rest"
)

type Options struct {
	IsDev           bool
	Addr            string
	DatabaseConnStr string
	AccrualHost     string
}

func New(opts Options) *App {
	server := rest.NewServer(rest.ServerOptions{
		IsDev: opts.IsDev,
		Addr:  opts.Addr,
	})

	return &App{
		server: server,
	}
}

type App struct {
	server *rest.Server
}

func (app *App) Run(ctx context.Context) {
	app.server.Run(ctx)
}

func (app *App) Wait() {
	app.server.Wait()
}
