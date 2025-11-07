package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"courier-service/internal/app"
	"courier-service/internal/core"
	"courier-service/internal/handlers"
	"courier-service/internal/repository"
	"courier-service/internal/usecase"
)

func main() {	
	cfg, _ := core.LoadConfig()

	dbPool := core.InitPool(context.Background())
	courierRepository := repository.NewCourierRepository(dbPool)
	courierUseCase := usecase.NewCourierUseCase(courierRepository)
	courierController := handlers.NewCourierController(courierUseCase)

	app := app.New(cfg.Port, courierController)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
