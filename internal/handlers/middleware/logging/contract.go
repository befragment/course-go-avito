package middleware

import "net/http"

type pathNormalizer interface {
	Normalize(r *http.Request) string
}

type metricsWriter interface {
	RecordRequest(method, path, status string)
	RecordDuration(method, path, status string, duration float64)
}
