package middleware

import "net/http"

type pathNormalizer interface {
	Normalize(r *http.Request) string
}

type metricsWriter interface {
	RecordRateLimitExceeded(method, path string)
}

type rateLimiter interface {
	Allow() bool
}
