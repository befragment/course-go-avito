package routing

import (
	"courier-service/internal/handlers/common"

	"github.com/go-chi/chi/v5"
)

func registerCommonRoutes(r chi.Router) {
	r.Get("/ping", common.Ping)
	r.Get("/health", common.Healthcheck)
}
