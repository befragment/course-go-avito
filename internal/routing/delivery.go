package routing

import (
	"github.com/go-chi/chi/v5"
)

func registerDeliveryRoutes(r chi.Router, c deliveryHandler) {
	r.Post("/delivery/assign", c.AssignDelivery)
	r.Post("/delivery/unassign", c.UnassignDelivery)
}
