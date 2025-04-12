package app

import (
	"context"

	"github.com/xantinium/gophermart/internal/infrastructure/memstorage"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres"
	"github.com/xantinium/gophermart/internal/presentation/rest"
	"github.com/xantinium/gophermart/internal/repository/tokensrepo"
	usersrepo "github.com/xantinium/gophermart/internal/repository/users"
	"github.com/xantinium/gophermart/internal/usecases"
)

type Options struct {
	IsDev           bool
	Addr            string
	DatabaseConnStr string
	AccrualHost     string
}

func New(opts Options) *App {
	psqlClient, err := postgres.NewPostgresClient(context.Background(), opts.DatabaseConnStr, postgres.DefaultOptions)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(rest.ServerOptions{
		IsDev: opts.IsDev,
		Addr:  opts.Addr,
		UseCases: usecases.NewUseCases(usecases.UseCasesOptions{
			UsersRepo:  usersrepo.NewUsersRepository(psqlClient),
			TokensRepo: tokensrepo.NewTokensRepository(memstorage.New()),
		}),
	})

	return &App{
		server:     server,
		psqlClient: psqlClient,
	}
}

type App struct {
	server     *rest.Server
	psqlClient *postgres.PostgresClient
}

func (app *App) Run(ctx context.Context) {
	app.server.Run(ctx)
}

func (app *App) Wait() {
	app.server.Wait()
	app.psqlClient.Destroy()
}
