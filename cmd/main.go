package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"courier-service/internal/app"
	"courier-service/internal/core"
)

func main() {
	cfg, _ := core.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.New(cfg.Port)

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
