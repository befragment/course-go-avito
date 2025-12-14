package main

import (
	"context"
	"log"

	core "courier-service/internal/core"
	ordergw "courier-service/internal/gateway/order"
	handlers "courier-service/internal/handlers"
	repository "courier-service/internal/repository"
	routing "courier-service/internal/routing"
	usecase "courier-service/internal/usecase"
	orderpb "courier-service/proto/order"

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

	dbPool := core.MustInitPool()
	// TODO: add grpc client to config
	grpcClient, _ := grpc.NewClient("service-order:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	ordersClient := orderpb.NewOrdersServiceClient(grpcClient)
	orderGateway := ordergw.NewGateway(ordersClient)
	courierRepo := repository.NewCourierRepository(dbPool)
	deliveryRepo := repository.NewDeliveryRepository(dbPool)
	txRunner := repository.NewTxRunner(dbPool)

	deliveryCalculator := usecase.NewFactory()
	deliveryUseCase := usecase.NewDelieveryUseCase(courierRepo, deliveryRepo, txRunner, deliveryCalculator)
	courierUseCase := usecase.NewCourierUseCase(courierRepo, deliveryCalculator)
	orderUseCase := usecase.NewOrderUsecase(orderGateway, courierRepo, deliveryRepo, txRunner, deliveryCalculator)

	go courierUseCase.CheckFreeCouriersWithInterval(backgroundCtx, cfg.CheckFreeCouriersInterval)
	go orderUseCase.ProcessOrders(backgroundCtx, cfg.OrderCheckCursorDelta)

	core.StartServer(
		dbPool,
		grpcClient,
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
