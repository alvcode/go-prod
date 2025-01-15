package main

import (
	"context"
	"prod/internal/app"
	"prod/internal/config"
	"prod/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.NewLogger()
	ctx = logging.ContextWithLogger(ctx, logger)

	logger.Infoln("Starting application")
	cfg := config.GetConfig()

	logger.Println("Loading config")

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}

	logging.GetLogger(ctx).Println("Before Run")
	a.Run(ctx)
}
