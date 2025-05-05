package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/xantinium/gophermart/internal/app"
	"github.com/xantinium/gophermart/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := parseAppArgs()

	logger.Init(args.IsDev)
	defer logger.Destroy()

	app := app.New(app.Options{
		IsDev:           args.IsDev,
		Addr:            args.Addr,
		DatabaseConnStr: args.DatabaseConnStr,
		AccrualHost:     args.AccrualHost,
	})

	app.Run(ctx)

	waitForStopSignal()
	cancel()

	app.Wait()
}

func waitForStopSignal() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan
}
