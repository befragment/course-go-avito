package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	core "courier-service/internal/core"
	interceptor "courier-service/internal/gateway/interceptor"
	ordergw "courier-service/internal/gateway/order"
	retryexec "courier-service/internal/gateway/retry"
	orderhandler "courier-service/internal/handlers/queues/order/changed"
	model "courier-service/internal/model"
	courierRepo "courier-service/internal/repository/courier"
	deliveryRepo "courier-service/internal/repository/delivery"
	txRunner "courier-service/internal/repository/txrunner"
	deliveryassignusecase "courier-service/internal/usecase/delivery/assign"
	deliverycompleteusecase "courier-service/internal/usecase/delivery/complete"
	deliveryunassignusecase "courier-service/internal/usecase/delivery/unassign"
	changed "courier-service/internal/usecase/order/changed"
	processor "courier-service/internal/usecase/order/changed/processor"
	deliverycalculator "courier-service/internal/usecase/utils"
	l "courier-service/pkg/logger"
	metrics "courier-service/pkg/metrics"
	shutdown "courier-service/pkg/shutdown"
	orderpb "courier-service/proto/order"
)

func main() {
	ctx := shutdown.WaitForShutdown()

	logger, err := l.New(l.LogLevelInfo)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	go initMetricsServer(ctx, ":9101", logger)

	config := sarama.NewConfig()
	configureKafkaClient(config)
	logger.Info("Kafka client configured")

	topic := cfg.KafkaTopic
	groupID := cfg.KafkaGroupID
	brokers := cfg.KafkaBrokers

	// Создаем метрики для worker
	httpMetrics := metrics.NewHTTPMetrics(prometheus.DefaultRegisterer)
	metricsWriter := metrics.NewMetricsWriter(httpMetrics)

	grpcClient, err := grpc.NewClient(
		cfg.GRPCServiceOrderServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.LoggingMetricsInterceptor(logger, metricsWriter),
		),
	)
	if err != nil {
		logger.Errorf("Failed to create grpc client: %v", err)
	}
	defer func() {
		if err := grpcClient.Close(); err != nil {
			logger.Errorf("Failed to close grpc client: %v", err)
		}
	}()

	ordersClient := orderpb.NewOrdersServiceClient(grpcClient)
	retryCfg := configureRetry(cfg.RetryMaxAttempts)
	retry := retryexec.NewRetryExecutor(retryCfg, logger)
	orderGateway := ordergw.NewGateway(ordersClient, retry, logger)

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
	logger *l.Logger,
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

func configureRetry(maxAttemps int) retryexec.RetryConfig {
	fullJitter := retryexec.NewFullJitter(50*time.Millisecond, 1*time.Second, 2.0)
	return retryexec.RetryConfig{
		MaxAttempts: maxAttemps,
		Strategy:    fullJitter,
		ShouldRetry: nil,
	}
}

func initMetricsServer(ctx context.Context, addr string, logger *l.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("metrics server listen failed on %s: %v", addr, err)
	}
}
