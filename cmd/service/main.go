package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	logger "courier-service/pkg/logger"
	core "courier-service/internal/core"

	courierRepo "courier-service/internal/repository/courier"
	deliveryRepo "courier-service/internal/repository/delivery"
	txRunner "courier-service/internal/repository/txrunner"

	courierusecase "courier-service/internal/usecase/courier"
	deliveryassignusecase "courier-service/internal/usecase/delivery/assign"
	deliveryunassignusecase "courier-service/internal/usecase/delivery/unassign"
	deliverycalculator "courier-service/internal/usecase/utils"

	courierhandlers "courier-service/internal/handlers/courier"
	deliveryhandlers "courier-service/internal/handlers/delivery"

	routing "courier-service/internal/routing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := core.WaitForShutdown()

	logger, err := logger.New(logger.LogLevelInfo)
	if err != nil { log.Fatalf("Failed to create logger: %v", err) }

	cfg, err := core.LoadConfig()
	if err != nil { logger.Error("Failed to load config: %v", err) }

	dbPool := core.MustInitPool(logger)
	defer dbPool.Close()

	grpcClient, err := grpc.NewClient(cfg.GRPCServiceOrderServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { logger.Errorf("Failed to create grpc client: %v", err) }
	defer grpcClient.Close()
	courierRepo := courierRepo.NewCourierRepository(dbPool, logger)
	deliveryRepo := deliveryRepo.NewDeliveryRepository(dbPool)
	txRunner := txRunner.NewTxRunner(dbPool)

	deliveryCalculator := deliverycalculator.NewTimeCalculatorFactory()
	assignUseCase := deliveryassignusecase.NewAssignDelieveryUseCase(
		courierRepo, 
		deliveryRepo, 
		txRunner, 
		deliveryCalculator,
	)
	unassignUseCase := deliveryunassignusecase.NewUnassignDelieveryUseCase(
		courierRepo, 
		deliveryRepo, 
		txRunner,
	)
	courierUseCase := courierusecase.NewCourierUseCase(
		courierRepo, 
		deliveryCalculator,
		logger,
	)

	go courierUseCase.CheckFreeCouriersWithInterval(ctx, cfg.CheckFreeCouriersInterval)

	router := routing.Router(
		logger,
		courierhandlers.NewCourierController(
			courierUseCase,
		),
		deliveryhandlers.NewDeliveryController(
			assignUseCase,
			unassignUseCase,
		),
	)

	startServer(ctx, cfg.Port, router, logger)
	<-ctx.Done()
	logger.Info("Service stopped gracefully")
}


func startServer(ctx context.Context, port string, handler http.Handler, logger logger.Interface) {
    srv := &http.Server{
        Addr:    port,
        Handler: handler,
    }

    serverErr := make(chan error, 1)
    go func() {
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            serverErr <- err
        }
        close(serverErr)
    }()

    select {
    case <-ctx.Done():
        logger.Info("Shutdown signal received")
    case err := <-serverErr:
        if err != nil {
            logger.Errorf("Server error: %v", err)
        }
    }

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("error shutting down server: %v", err)
	}
}
