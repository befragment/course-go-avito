package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingMetricsInterceptor создает unary interceptor для логирования и метрик
func LoggingMetricsInterceptor(logger logger, metrics metricsWriter) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		logger.Infof("gRPC request started: method=%s", method)

		// Выполняем запрос
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Измеряем длительность
		duration := time.Since(start)

		// Логируем результат
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			logger.Errorf("gRPC request failed: method=%s, duration=%v, code=%s, error=%v",
				method, duration, grpcStatus.Code(), err)

			// Записываем метрику для каждой retryable ошибки
			// Для gRPC: HTTP метод всегда "POST", путь = gRPC метод
			if isRetryableError(grpcStatus.Code()) {
				metrics.RecordRetry("POST", method)
			}
		} else {
			logger.Infof("gRPC request succeeded: method=%s, duration=%v",
				method, duration)
		}

		return err
	}
}

// MetricsInterceptor создает unary interceptor только для метрик
func MetricsInterceptor(metrics metricsWriter) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Записываем метрику если была ошибка
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			// Для gRPC: HTTP метод всегда "POST", путь = gRPC метод
			if isRetryableError(grpcStatus.Code()) {
				metrics.RecordRetry("POST", method)
			}
		}

		return err
	}
}

// LoggingInterceptor создает unary interceptor только для логирования (без метрик)
func LoggingInterceptor(logger logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		// Логируем начало запроса
		logger.Infof("gRPC request started: method=%s", method)

		// Выполняем запрос
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Измеряем длительность
		duration := time.Since(start)

		// Логируем результат
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			logger.Errorf("gRPC request failed: method=%s, duration=%v, code=%s, error=%v",
				method, duration, grpcStatus.Code(), err)
		} else {
			logger.Infof("gRPC request succeeded: method=%s, duration=%v",
				method, duration)
		}

		return err
	}
}

// isRetryableError определяет, является ли ошибка повторяемой
func isRetryableError(code codes.Code) bool {
	switch code {
	case codes.Unavailable,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal:
		return true
	default:
		return false
	}
}
