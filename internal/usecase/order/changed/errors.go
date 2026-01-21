package changed

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderStatusMismatch = errors.New("order status mismatch")
)
