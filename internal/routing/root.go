package routing

import (
	"github.com/go-chi/chi/v5"

	loggingmiddleware "courier-service/internal/handlers/middleware/logging"
	ratelimitmiddleware "courier-service/internal/handlers/middleware/ratelimit"
)

func Router(
	logger logger,
	rateLimiter rateLimiter,
	metricsWriter httpMetricsWriter,
	metricsHandler metricsHandler,
	pathNormalizer pathNormalizer,
	courierController courierHandler,
	deliveryController deliveryHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// /metrics endpoint БЕЗ rate limiting
	r.Handle("/metrics", metricsHandler)

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
