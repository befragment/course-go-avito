package routing

import (
	"github.com/go-chi/chi/v5"

	"courier-service/internal/handlers"
)

func registerCourierRoutes(r chi.Router, c *handlers.CourierController) {
	r.Get("/couriers", c.GetAll)
	r.Route("/courier", func(r chi.Router) {
		r.Get("/{id}", c.GetById)
		r.Post("/", c.Create)
		r.Put("/", c.Update)
	})
}
