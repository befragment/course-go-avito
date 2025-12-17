package courier

import "errors"

var (
	ErrCourierNotFound = errors.New("courier not found")
	ErrCouriersBusy = errors.New("all couriers are busy")
	ErrOrderIDExists = errors.New("order id already exists")
	ErrOrderIDNotFound = errors.New("order id not found")
	ErrPhoneNumberExists = errors.New("phone number already exists")
	ErrNothingToUpdate   = errors.New("nothing to update")
	ErrOrderNotFound     = errors.New("order not found")
)