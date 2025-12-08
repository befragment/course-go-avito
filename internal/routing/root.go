package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"courier-service/internal/handlers"
)
func Router(
	courierController *handlers.CourierController,
	deliveryController *handlers.DeliveryController,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	registerCourierRoutes(r, courierController)
	registerDeliveryRoutes(r, deliveryController)
	return r
}