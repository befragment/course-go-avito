package routing

import "net/http"

type httpMetricsWriter interface {
	RecordRequest(method, path, status string)
	RecordDuration(method, path, status string, duration float64)
	RecordRetry(method, path string)
	RecordRateLimitExceeded(method, path string)
}

type pathNormalizer interface {
	Normalize(r *http.Request) string
}

type courierHandler interface {
	GetCourierById(w http.ResponseWriter, r *http.Request)
	GetAllCouriers(w http.ResponseWriter, r *http.Request)
	CreateCourier(w http.ResponseWriter, r *http.Request)
	UpdateCourier(w http.ResponseWriter, r *http.Request)
}

type deliveryHandler interface {
	AssignDelivery(w http.ResponseWriter, r *http.Request)
	UnassignDelivery(w http.ResponseWriter, r *http.Request)
}

type metricsHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type rateLimiter interface {
	Allow() bool
}
