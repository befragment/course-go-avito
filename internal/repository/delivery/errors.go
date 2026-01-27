package delivery

import "errors"

var (
	ErrOrderIDExists   = errors.New("order id already exists")
	ErrOrderIDNotFound = errors.New("order id not found")
)
