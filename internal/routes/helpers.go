package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterPing(r chi.Router) {
	r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})
}

func RegisterHealthcheck(r chi.Router) {
	r.Head("/healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	
	RegisterPing(r)
	RegisterHealthcheck(r)

	RegisterCouriersRoutes(r)
	return r
}