package interceptor

type metricsWriter interface {
	RecordRetry(method, path string)
}