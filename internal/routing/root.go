package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"courier-service/internal/handlers/courier"
	"courier-service/internal/handlers/delivery"
)
func Router(
	courierController *courier.CourierController,
	deliveryController *delivery.DeliveryController,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	registerCourierRoutes(r, courierController)
	registerDeliveryRoutes(r, deliveryController)
	return r
}