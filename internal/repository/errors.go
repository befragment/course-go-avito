package repository

import "errors"

var (
	ErrPhoneNumberExists = errors.New("phone number already exists")
	ErrNothingToUpdate   = errors.New("nothing to update")
	ErrCourierNotFound   = errors.New("courier not found")
)