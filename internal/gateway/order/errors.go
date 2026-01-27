package order

import "errors"

var (
	ErrRetryLimitExceeded = errors.New("retry limit exceeded")
)
