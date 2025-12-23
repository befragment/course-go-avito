package routing

import (
	"github.com/go-chi/chi/v5"

	"courier-service/internal/handlers/courier"
)

func registerCourierRoutes(r chi.Router, c *courier.CourierController) {
	r.Get("/couriers", c.GetAllCouriers)
	r.Route("/courier", func(r chi.Router) {
		r.Get("/{id}", c.GetCourierById)
		r.Post("/", c.CreateCourier)
		r.Put("/", c.UpdateCourier)
	})
}
