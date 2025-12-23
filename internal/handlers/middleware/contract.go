package middleware

import "net/http"

type pathNormalizer interface {
	Normalize(r *http.Request) string
}

