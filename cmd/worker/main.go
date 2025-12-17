package main

import (
	"context"
	"log"
	"time"
	"github.com/IBM/sarama"
	logger "courier-service/pkg/logger"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"

	core "courier-service/internal/core"
	// orderpb "courier-service/proto/order"
	// ordergw "courier-service/internal/gateway/order"
	courierRepo "courier-service/internal/repository/courier"
	deliveryRepo "courier-service/internal/repository/delivery"
	txRunner "courier-service/internal/repository/utils/txrunner"
	orderhandler "courier-service/internal/transport/kafka/order"

	deliveryassignusecase "courier-service/internal/usecase/delivery/assign"
	deliveryunassignusecase "courier-service/internal/usecase/delivery/unassign"
	deliverycompleteusecase "courier-service/internal/usecase/delivery/complete"
	deliverycalculator "courier-service/internal/usecase/utils"
)

func main() {
	backgroundCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("kafka worker started")

	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := logger.New(logger.LogLevelInfo)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	config := sarama.NewConfig()
	configureKafkaClient(config)

	topic := cfg.KafkaTopic
	groupID := cfg.KafkaGroupID
	brokers := cfg.KafkaBrokers

	dbPool := core.MustInitPool()
	courierRepository := courierRepo.NewCourierRepository(dbPool, logger)
	deliveryRepository := deliveryRepo.NewDeliveryRepository(dbPool)
	transactionRunner := txRunner.NewTxRunner(dbPool)
	// grpcClient, _ := grpc.NewClient(cfg.GRPCServiceOrderServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// ordersClient := orderpb.NewOrdersServiceClient(grpcClient)
	// orderGateway := ordergw.NewGateway(ordersClient)

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
	
	orderChangedFactory := orderhandler.NewOrderChangedFactory(
		assignUseCase,
		unassignUseCase,
		completeUseCase,
	)
	
	orderChangeHandler := orderhandler.NewOrderStatusChangedHandler(
		orderChangedFactory,
	)

	if err := runKafkaConsumer(backgroundCtx, brokers, groupID, topic, config, orderChangeHandler); err != nil {
		log.Fatalf("Kafka consumer stopped with error: %v", err)
	}

	log.Println("Kafka consumer exited gracefully")
}

func configureKafkaClient(config *sarama.Config) {
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
}

func runKafkaConsumer(
	ctx context.Context,
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
			log.Printf("Failed to close kafka client: %v", err)
		}
	}()

	for {
		if err := client.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("Error from consumer: %v", err)
			time.Sleep(time.Second)
		}

		if ctx.Err() != nil {
			log.Println("Context cancelled, exiting consume loop")
			return nil
		}
	}
}