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
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return "unknown"
	}

	// If router already matched (e.g. after next.ServeHTTP)
	if rp := rctx.RoutePattern(); rp != "" {
		return rp
	}

	// Early stage: do a dry-run match against the same Routes tree
	if rctx.Routes != nil {
		tmp := chi.NewRouteContext()
		if rctx.Routes.Match(tmp, r.Method, r.URL.Path) {
			if rp := tmp.RoutePattern(); rp != "" {
				return rp
			}
		}
	}

	return "unknown"
}
