package main

import (
	"context"
	"log"

	logger "courier-service/pkg/logger"
	core "courier-service/internal/core"

	courierRepo "courier-service/internal/repository/courier"
	deliveryRepo "courier-service/internal/repository/delivery"
	txRunner "courier-service/internal/repository/utils/txrunner"
	// ordergw "courier-service/internal/gateway/order"

	courierusecase "courier-service/internal/usecase/courier"
	deliveryassignusecase "courier-service/internal/usecase/delivery/assign"
	deliveryunassignusecase "courier-service/internal/usecase/delivery/unassign"
	// orderusecase "courier-service/internal/usecase/order"
	deliverycalculator "courier-service/internal/usecase/utils"

	courierhandlers "courier-service/internal/handlers/courier"
	deliveryhandlers "courier-service/internal/handlers/delivery"

	routing "courier-service/internal/routing"

	// orderpb "courier-service/proto/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	backgroundCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := logger.New(logger.LogLevelInfo)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	dbPool := core.MustInitPool()

	grpcClient, _ := grpc.NewClient(cfg.GRPCServiceOrderServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// ordersClient := orderpb.NewOrdersServiceClient(grpcClient)
	// orderGateway := ordergw.NewGateway(ordersClient)
	courierRepo := courierRepo.NewCourierRepository(dbPool, logger)
	deliveryRepo := deliveryRepo.NewDeliveryRepository(dbPool)
	txRunner := txRunner.NewTxRunner(dbPool)

	deliveryCalculator := deliverycalculator.NewTimeCalculatorFactory()
	assignUseCase := deliveryassignusecase.NewAssignDelieveryUseCase(courierRepo, deliveryRepo, txRunner, deliveryCalculator)
	unassignUseCase := deliveryunassignusecase.NewUnassignDelieveryUseCase(courierRepo, deliveryRepo, txRunner)
	courierUseCase := courierusecase.NewCourierUseCase(courierRepo, deliveryCalculator)
	// orderMonitoringUseCase := orderusecase.NewOrderMonitoringUseCase(orderGateway, courierRepo, deliveryRepo, txRunner, deliveryCalculator)

	go courierUseCase.CheckFreeCouriersWithInterval(backgroundCtx, cfg.CheckFreeCouriersInterval)
	// go orderMonitoringUseCase.MonitorOrders(backgroundCtx, cfg.OrderCheckCursorDelta)

	core.StartServer(
		dbPool,
		grpcClient,
		cfg.Port,
		routing.Router(
			courierhandlers.NewCourierController(
				courierUseCase,
			),
			deliveryhandlers.NewDeliveryController(
				assignUseCase,
				unassignUseCase,
			),
		),
	)
}
