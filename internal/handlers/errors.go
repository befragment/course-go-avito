package handlers

import (
	"log"
	"net/http"

	"courier-service/internal/usecase"
)

const (
	ErrInvalidID             = "invalid id"
	ErrCourierNotFound       = "courier not found"
	ErrMissingRequiredFields = "Missing required fields"
	ErrInvalidPhoneNumber    = "Invalid phone number"
	ErrPhoneAlreadyExists    = "Phone number already exists"
	ErrInternalServer        = "Internal server error"
	ErrIDRequired            = "id is required"
)

func handleCreateError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrInvalidCreate:
		respondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case usecase.ErrInvalidPhoneNumber:
		respondWithError(w, http.StatusBadRequest, ErrInvalidPhoneNumber)
	case usecase.ErrPhoneNumberExists:
		respondWithError(w, http.StatusConflict, ErrPhoneAlreadyExists)
	default:
		log.Printf("Internal server error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
	}
}

func handleUpdateError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrInvalidUpdate:
		respondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case usecase.ErrInvalidPhoneNumber:
		respondWithError(w, http.StatusBadRequest, ErrInvalidPhoneNumber)
	case usecase.ErrPhoneNumberExists:
		respondWithError(w, http.StatusConflict, ErrPhoneAlreadyExists)
	case usecase.ErrCourierNotFound:
		respondWithError(w, http.StatusNotFound, ErrCourierNotFound)
	default:
		log.Printf("Internal server error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
	}
}
