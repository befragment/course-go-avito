package routing

import (
	"github.com/go-chi/chi/v5"

	loggingmiddleware "courier-service/internal/handlers/middleware/logging"
	ratelimitmiddleware "courier-service/internal/handlers/middleware/ratelimit"
	logger "courier-service/pkg/logger"
	ratelimiter "courier-service/pkg/ratelimiter"
)

func Router(
	logger logger.LoggerInterface,
	rateLimiter ratelimiter.RateLimiterInterface,
	metricsWriter httpMetricsWriter,
	metricsHandler metricsHandler,
	pathNormalizer pathNormalizer,
	courierController courierHandler,
	deliveryController deliveryHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// /metrics endpoint БЕЗ rate limiting (для Prometheus)
	r.Handle("/metrics", metricsHandler)

	// Группа маршрутов С middleware
	r.Group(func(r chi.Router) {
		r.Use(
			ratelimitmiddleware.RateLimitMiddleware(
				rateLimiter,
				logger,
				metricsWriter,
				pathNormalizer,
			),
			loggingmiddleware.LoggingMiddleware(
				logger,
				metricsWriter,
				pathNormalizer,
			),
		)

		registerCommonRoutes(r)
		registerCourierRoutes(r, courierController)
		registerDeliveryRoutes(r, deliveryController)
	})

	return r
}
