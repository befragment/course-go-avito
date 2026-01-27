package unassign

import "errors"

var (
	ErrPhoneNumberExists  = errors.New("phone number already exists")
	ErrNothingToUpdate    = errors.New("nothing to update")
	ErrCourierNotFound    = errors.New("courier not found")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidCreate      = errors.New("name, phone and status are required")
	ErrIdRequired         = errors.New("id is required")
	ErrInvalidUpdate      = errors.New("no fields provided for update")
	ErrCouriersBusy       = errors.New("all couriers are busy")

	ErrUnknownTransportType = errors.New("unknown transport type")
	ErrNoOrderID            = errors.New("order id is required")
	ErrOrderIDExists        = errors.New("order id already exists")
	ErrOrderIDNotFound      = errors.New("order id not found")
)
