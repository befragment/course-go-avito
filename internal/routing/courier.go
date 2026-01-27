package routing

import (
	"github.com/go-chi/chi/v5"
)

func registerCourierRoutes(r chi.Router, c courierHandler) {
	r.Get("/couriers", c.GetAllCouriers)
	r.Get("/courier/{id}", c.GetCourierById)
	r.Post("/courier", c.CreateCourier)
	r.Put("/courier", c.UpdateCourier)
}
