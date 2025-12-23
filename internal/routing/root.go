package routing

import (
	"github.com/go-chi/chi/v5"
	"courier-service/internal/handlers/courier"
	"courier-service/internal/handlers/delivery"
	"courier-service/internal/handlers/metrics"
	middleware "courier-service/internal/handlers/middleware"
	logger "courier-service/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(
	logger logger.Interface,
	courierController *courier.CourierController,
	deliveryController *delivery.DeliveryController,
) *chi.Mux {
	r := chi.NewRouter()
	normalizer := NewChiPathNormalizer()
	httpMetrics := metrics.NewHTTPMetrics(prometheus.DefaultRegisterer)

	r.Use(
		middleware.LoggingMiddleware(
			logger,
			normalizer,
			httpMetrics,
		),
	)
	
	registerCommonRoutes(r)
	registerCourierRoutes(r, courierController)
	registerDeliveryRoutes(r, deliveryController)

	r.Handle("/metrics", promhttp.Handler())
	return r
}
