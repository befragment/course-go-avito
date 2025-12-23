package middleware

import (
	"net/http"
	"strconv"
	"time"

	promMetrics "courier-service/internal/handlers/metrics"
	loggerpkg "courier-service/pkg/logger"
)



func LoggingMiddleware(logger loggerpkg.Interface, normalizer pathNormalizer, m *promMetrics.HTTPMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			rr := &responseRecorder{ResponseWriter: w, status: 200}
			next.ServeHTTP(rr, r)
			path := normalizer.Normalize(r)
			
			ignored := map[string]bool{
				"/metrics": true,
				"/health":  true,
			}
			if ignored[path] {
				logger.Debugf("skipping path: %s", path)
				return
			}
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(rr.status)

			m.RequestTotal.WithLabelValues(r.Method, path, status).Inc()
			m.RequestDuration.WithLabelValues(path, r.Method).Observe(duration)

			logger.Infof(loggerpkg.PrettyRequestLogFormat,
				time.Now().Format("2006/01/02 15:04:05"),
				r.Method,
				path,
				rr.status,
				duration*1000,
			)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
