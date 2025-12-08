package repository

import "errors"

var (
	ErrPhoneNumberExists = errors.New("phone number already exists")
	ErrNothingToUpdate   = errors.New("nothing to update")
	ErrCourierNotFound   = errors.New("courier not found")
	ErrCouriersBusy      = errors.New("all couriers are busy")
	ErrOrderIDExists     = errors.New("order id already exists")
	ErrOrderIDNotFound   = errors.New("order id not found")
)