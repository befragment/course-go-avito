package main

import (
	"context"
	logger "courier-service/pkg/logger"
	orderpb "courier-service/proto/order"
	"log"
	"time"

	"github.com/IBM/sarama"

	core "courier-service/internal/core"
	ordergw "courier-service/internal/gateway/order"
	orderhandler "courier-service/internal/handlers/queues/order/changed"
	courierRepo "courier-service/internal/repository/courier"
	deliveryRepo "courier-service/internal/repository/delivery"
	txRunner "courier-service/internal/repository/txrunner"

	model "courier-service/internal/model"
	deliveryassignusecase "courier-service/internal/usecase/delivery/assign"
	deliverycompleteusecase "courier-service/internal/usecase/delivery/complete"
	deliveryunassignusecase "courier-service/internal/usecase/delivery/unassign"
	changed "courier-service/internal/usecase/order/changed"
	processor "courier-service/internal/usecase/order/changed/processor"
	deliverycalculator "courier-service/internal/usecase/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := core.WaitForShutdown()

	logger, err := logger.New(logger.LogLevelInfo)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
		
	config := sarama.NewConfig()
	configureKafkaClient(config)
	logger.Info("Kafka client configured")

	topic := cfg.KafkaTopic
	groupID := cfg.KafkaGroupID
	brokers := cfg.KafkaBrokers

	grpcClient, err := grpc.NewClient(cfg.GRPCServiceOrderServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Failed to create grpc client: %v", err)
	}
	defer grpcClient.Close()

	ordersClient := orderpb.NewOrdersServiceClient(grpcClient)
	orderGateway := ordergw.NewGateway(ordersClient)

	dbPool := core.MustInitPool(logger)
	defer dbPool.Close()
	courierRepository := courierRepo.NewCourierRepository(dbPool, logger)
	deliveryRepository := deliveryRepo.NewDeliveryRepository(dbPool)
	transactionRunner := txRunner.NewTxRunner(dbPool)

	deliveryCalculator := deliverycalculator.NewTimeCalculatorFactory()

	assignUseCase := deliveryassignusecase.NewAssignDelieveryUseCase(
		courierRepository,
		deliveryRepository,
		transactionRunner,
		deliveryCalculator,
	)
	unassignUseCase := deliveryunassignusecase.NewUnassignDelieveryUseCase(
		courierRepository,
		deliveryRepository,
		transactionRunner,
	)
	completeUseCase := deliverycompleteusecase.NewCompleteDeliveryUseCase(
		courierRepository,
	)

	createdProcessor := processor.NewCreatedProcessor(assignUseCase)
	cancelledProcessor := processor.NewCancelledProcessor(unassignUseCase)
	completedProcessor := processor.NewCompletedProcessor(completeUseCase)

	orderChangedFactory := changed.NewFactory(map[model.OrderStatus]changed.Processor{
		model.OrderStatusCreated:   createdProcessor,
		model.OrderStatusCancelled: cancelledProcessor,
		model.OrderStatusCompleted: completedProcessor,
	})

	orderChangedUseCase := changed.NewOrderChangedUseCase(orderChangedFactory, orderGateway, logger)
	orderChangedHandler := orderhandler.NewOrderStatusChangedHandler(orderChangedUseCase, logger)

	go func() {
		if err := runKafkaConsumer(ctx, logger, brokers, groupID, topic, config, orderChangedHandler); err != nil {
			logger.Errorf("Kafka consumer stopped with error: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Kafka consumer exited gracefully")
}

func configureKafkaClient(config *sarama.Config) {
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
}

func runKafkaConsumer(
	ctx context.Context,
	logger logger.Interface,
	brokers []string,
	groupID string,
	topic string,
	config *sarama.Config,
	handler sarama.ConsumerGroupHandler,
) error {
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			logger.Errorf("Failed to close kafka client: %v", err)
		}
	}()

	for {
		if err := client.Consume(ctx, []string{topic}, handler); err != nil {
			logger.Errorf("Error from consumer: %v", err)
			time.Sleep(time.Second)
		}

		if ctx.Err() != nil {
			logger.Info("Context cancelled, exiting consume loop")
			return nil
		}
	}
}
