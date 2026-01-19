package middleware

import (
	"net/http"

	lpkg "courier-service/pkg/logger"
	rlpkg "courier-service/pkg/ratelimiter"
)

func RateLimitMiddleware(
	limiter rlpkg.RateLimiterInterface,
	logger lpkg.LoggerInterface,
	metricsWriter metricsWriter,
	normalizer pathNormalizer,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				path := normalizer.Normalize(r)
				logger.Warnf("Rate limit exceeded for %s", path)
				metricsWriter.RecordRateLimitExceeded(r.Method, path)
				w.Header().Set("X-RateLimit-Limit", "10")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Rate limit exceeded"))

				return
			}

			logger.Debugf("Request allowed: %s", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
