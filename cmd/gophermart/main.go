package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/xantinium/gophermart/internal/app"
	"github.com/xantinium/gophermart/internal/logger"
)

const isDev = true

func main() {
	args := parseAppArgs()

	logger.Init(isDev)
	defer logger.Destroy()

	app := app.New(app.Options{
		IsDev:           isDev,
		Addr:            args.Addr,
		DatabaseConnStr: args.DatabaseConnStr,
		AccrualHost:     args.AccrualHost,
	})

	app.Run(context.Background())

	waitForStopSignal()

	app.Wait()
}

func waitForStopSignal() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan
}
