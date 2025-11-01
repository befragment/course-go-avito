package main

import (
	"context"
	"os/signal"
	"syscall"
	"log"
	"os"

	"github.com/Avito-courses/course-go-avito-befragment/internal/app"
	"github.com/Avito-courses/course-go-avito-befragment/internal/core"
)


func main() {
	cfg, _ := core.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a := app.New(cfg.Port)

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}