package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics struct {
	RequestTotal           *prometheus.CounterVec
	RequestDuration        *prometheus.HistogramVec
	RateLimitExceededTotal *prometheus.CounterVec
	GatewayRetries         *prometheus.CounterVec
}

func NewHTTPMetrics(reg prometheus.Registerer) *HTTPMetrics {
	metrics := &HTTPMetrics{
		RequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		RateLimitExceededTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_exceeded_total",
				Help: "Total number of rate limiting",
			},
			[]string{"method", "path"},
		),
		GatewayRetries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_retries_total",
				Help: "Total number of gateway retries",
			},
			[]string{"method", "path"},
		),
	}
	reg.MustRegister(
		metrics.RequestTotal,
		metrics.RequestDuration,
		metrics.RateLimitExceededTotal,
		metrics.GatewayRetries,
	)
	return metrics
}

type MetricsWriter struct {
	metrics *HTTPMetrics
}

func NewMetricsWriter(metrics *HTTPMetrics) *MetricsWriter {
	return &MetricsWriter{
		metrics: metrics,
	}
}

func (w *MetricsWriter) RecordRequest(method, path, status string) {
	w.metrics.RequestTotal.WithLabelValues(method, path, status).Inc()
}

func (w *MetricsWriter) RecordDuration(method, path, status string, duration float64) {
	w.metrics.RequestDuration.WithLabelValues(method, path, status).Observe(duration)
}

func (w *MetricsWriter) RecordRetry(method, path string) {
	w.metrics.GatewayRetries.WithLabelValues(method, path).Inc()
}

func (w *MetricsWriter) RecordRateLimitExceeded(method, path string) {
	w.metrics.RateLimitExceededTotal.WithLabelValues(method, path).Inc()
}
