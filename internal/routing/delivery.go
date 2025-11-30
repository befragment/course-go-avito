package routing

import (
	"courier-service/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func registerDeliveryRoutes(r chi.Router, c *handlers.DeliveryController) {
	r.Route("/delivery", func(r chi.Router) {
		r.Post("/assign", c.AssignDelivery)
		r.Post("/unassign", c.UnassignDelivery)
	})
}