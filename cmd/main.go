package main

import (
	"context"
	"courier-service/internal/core"
	"courier-service/internal/handlers"
	"courier-service/internal/repository"
	"courier-service/internal/routing"
	"courier-service/internal/usecase"
	"log"
)

func main() {
	backgroundCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool := core.MustInitPool()

	courierRepo := repository.NewCourierRepository(dbPool)
	deliveryRepo := repository.NewDeliveryRepository(dbPool)
	txRunner := repository.NewTxRunner(dbPool)

	deliveryUseCase := usecase.NewDelieveryUseCase(courierRepo, deliveryRepo, txRunner)
	courierUseCase := usecase.NewCourierUseCase(courierRepo)

	core.TestDBConnString()

	go usecase.CheckFreeCouriers(backgroundCtx, courierUseCase)

	core.StartServer(
		dbPool,
		cfg.Port,
		routing.Router(
			handlers.NewCourierController(
				courierUseCase,
			),
			handlers.NewDeliveryController(
				deliveryUseCase,
			),
		),
	)
}
