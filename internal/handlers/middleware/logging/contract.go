package middleware

import "net/http"

type pathNormalizer interface {
	Normalize(r *http.Request) string
}

type metricsWriter interface {
	RecordRequest(method, path, status string)
	RecordDuration(method, path, status string, duration float64)
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
