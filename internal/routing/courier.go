package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"courier-service/internal/handlers"
)

func InitCourierRoutes(c *handlers.CourierController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", handlers.Ping)
	r.Head("/healthcheck", handlers.Healthcheck)

	r.Get("/couriers", c.GetAll)

	r.Route("/courier", func(r chi.Router) {
		r.Get("/{id}", c.GetById)
		r.Post("/", c.Create)
		r.Put("/", c.Update)
	})
	return r
}