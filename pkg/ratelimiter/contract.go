package ratelimiter

type RateLimiterInterface interface {
	Allow() bool
}
