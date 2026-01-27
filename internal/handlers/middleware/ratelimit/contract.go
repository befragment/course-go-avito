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

type logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
