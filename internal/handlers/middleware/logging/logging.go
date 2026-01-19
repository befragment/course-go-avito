package middleware

import (
	"net/http"
	"strconv"
	"time"

	l "courier-service/pkg/logger"
)

func LoggingMiddleware(
	logger l.LoggerInterface,
	metricsWriter metricsWriter,
	normalizer pathNormalizer,
) func(http.Handler) http.Handler {
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

			metricsWriter.RecordRequest(r.Method, path, status)
			metricsWriter.RecordDuration(r.Method, path, status, duration)

			logger.Infof(l.PrettyRequestLogFormat,
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
