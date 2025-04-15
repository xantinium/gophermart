package app

import (
	"context"

	"github.com/xantinium/gophermart/internal/app/worker"
	"github.com/xantinium/gophermart/internal/infrastructure/memstorage"
	"github.com/xantinium/gophermart/internal/infrastructure/postgres"
	"github.com/xantinium/gophermart/internal/presentation/rest"
	ordersrepo "github.com/xantinium/gophermart/internal/repository/orders"
	tokensrepo "github.com/xantinium/gophermart/internal/repository/tokens"
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

	cases := usecases.New(usecases.Options{
		UsersRepo:  usersrepo.New(psqlClient),
		TokensRepo: tokensrepo.New(memstorage.New()),
		OrdersRepo: ordersrepo.New(psqlClient),
	})

	workerPool := worker.NewWorkerPool(worker.WorkerPoolOptions{
		PoolSize:    10,
		AccrualHost: opts.AccrualHost,
		UseCases:    cases,
	})

	server := rest.NewServer(rest.ServerOptions{
		IsDev:    opts.IsDev,
		Addr:     opts.Addr,
		UseCases: cases,
	})

	return &App{
		workerPool: workerPool,
		server:     server,
		psqlClient: psqlClient,
	}
}

type App struct {
	workerPool *worker.WorkerPool
	server     *rest.Server
	psqlClient *postgres.PostgresClient
}

func (app *App) Run(ctx context.Context) {
	app.workerPool.Run(ctx)
	app.server.Run(ctx)
}

func (app *App) Wait() {
	app.workerPool.Wait()
	app.server.Wait()
	app.psqlClient.Destroy()
}
