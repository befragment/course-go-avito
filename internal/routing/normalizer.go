package routing

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ChiPathNormalizer struct{}

func NewChiPathNormalizer() *ChiPathNormalizer {
	return &ChiPathNormalizer{}
}

func (n *ChiPathNormalizer) Normalize(r *http.Request) string {
	routePattern := chi.RouteContext(r.Context()).RoutePattern()

	if routePattern == "" {
		return "unknown"
	}

	return routePattern
}
