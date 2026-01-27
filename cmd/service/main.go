package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	prometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	core "courier-service/internal/core"
	interceptor "courier-service/internal/gateway/interceptor"
	courierhandlers "courier-service/internal/handlers/courier"
	deliveryhandlers "courier-service/internal/handlers/delivery"
	courierRepo "courier-service/internal/repository/courier"
	deliveryRepo "courier-service/internal/repository/delivery"
	txRunner "courier-service/internal/repository/txrunner"
	routing "courier-service/internal/routing"
	courierusecase "courier-service/internal/usecase/courier"
	deliveryassignusecase "courier-service/internal/usecase/delivery/assign"
	deliveryunassignusecase "courier-service/internal/usecase/delivery/unassign"
	deliverycalculator "courier-service/internal/usecase/utils"
	database "courier-service/pkg/database/postgres"
	l "courier-service/pkg/logger/zap"
	metrics "courier-service/pkg/metrics/prometheus"
	rlimiter "courier-service/pkg/ratelimiter/tokenbucket"
	shutdown "courier-service/pkg/shutdown"
)

func main() {
	ctx := shutdown.WaitForShutdown()

	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger, err := l.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	ratelimiter := rlimiter.NewTokenBucket(cfg.TokenBucketCapacity, cfg.TokenBucketRefillRate, time.Now)

	dbPool := database.MustInitPool(cfg.PostgresDSN(), logger)
	defer dbPool.Close()

	// Создаем метрики до grpcClient, чтобы использовать в interceptor
	httpMetrics := metrics.NewHTTPMetrics(prometheus.DefaultRegisterer)
	metricsWriter := metrics.NewMetricsWriter(httpMetrics)

	// Создаем grpcClient с interceptors
	grpcClient, err := grpc.NewClient(
		cfg.GRPCServiceOrderServer,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.LoggingMetricsInterceptor(logger, metricsWriter),
		),
	)
	if err != nil {
		logger.Errorf("Failed to create grpc client: %v", err)
		return
	}

	defer func() {
		if err := grpcClient.Close(); err != nil {
			logger.Errorf("Failed to close grpc client: %v", err)
		}
	}()

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

	pathNormalizer := routing.NewChiPathNormalizer()
	metricsHandler := metrics.NewMetricsHandler()
	router := routing.Router(
		logger,
		ratelimiter,
		metricsWriter,
		metricsHandler,
		pathNormalizer,
		courierhandlers.NewCourierController(
			courierUseCase,
		),
		deliveryhandlers.NewDeliveryController(
			assignUseCase,
			unassignUseCase,
		),
	)
	logger.Info("Starting service server...")
	go startServer(ctx, cfg.Port, router, logger)
	logger.Info("Starting pprof server...")
	go startPprofServer(ctx, cfg.PprofAddress, logger)
	<-ctx.Done()
	logger.Info("Service stopped gracefully")
}

func startServer(ctx context.Context, port string, handler http.Handler, logger *l.Logger) {
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

func startPprofServer(ctx context.Context, addr string, logger *l.Logger) {
	logger.Infof("%s", addr)
	srv := &http.Server{
		Addr:              addr,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Infof("pprof listening on http://%s/debug/pprof/", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case <-ctx.Done():
		logger.Info("pprof shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			logger.Errorf("pprof server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("error shutting down pprof server: %v", err)
	}
}
