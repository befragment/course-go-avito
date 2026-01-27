package routing

import (
	"github.com/go-chi/chi/v5"

	"courier-service/internal/handlers/common"
)

func registerCommonRoutes(r chi.Router) {
	r.Get("/ping", common.Ping)
	r.Get("/health", common.Healthcheck)
}
