package handlers

import (
	"log"
	"net/http"

	"courier-service/internal/usecase"
)

const (
	ErrInvalidID             = "Invalid id"
	ErrCourierNotFound       = "Courier not found"
	ErrMissingRequiredFields = "Missing required fields"
	ErrInvalidPhoneNumber    = "Invalid phone number"
	ErrUnknownTransportType  = "Unknown transport type"
	ErrPhoneAlreadyExists    = "Phone number already exists"
	ErrInternalServer        = "Internal server error"
	ErrIDRequired            = "Id is required"
	ErrCouriersBusy          = "All couriers are busy"
	ErrOrderIDExists         = "Order id already exists"
	ErrNoCourierForOrder     = "No courier found for the order"
	ErrOrderIDNotFound       = "Order id not found"
)

func respondInternalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v\n", err)
	respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
}

func handleCreateError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrInvalidCreate:
		respondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case usecase.ErrInvalidPhoneNumber:
		respondWithError(w, http.StatusBadRequest, ErrInvalidPhoneNumber)
	case usecase.ErrUnknownTransportType:
		respondWithError(w, http.StatusBadRequest, ErrUnknownTransportType)
	case usecase.ErrPhoneNumberExists:
		respondWithError(w, http.StatusConflict, ErrPhoneAlreadyExists)
	default:
		respondInternalServerError(w, err)
	}
}

func handleUpdateError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrInvalidUpdate:
		respondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case usecase.ErrInvalidPhoneNumber:
		respondWithError(w, http.StatusBadRequest, ErrInvalidPhoneNumber)
	case usecase.ErrUnknownTransportType:
		respondWithError(w, http.StatusBadRequest, ErrUnknownTransportType)
	case usecase.ErrPhoneNumberExists:
		respondWithError(w, http.StatusConflict, ErrPhoneAlreadyExists)
	case usecase.ErrCourierNotFound:
		respondWithError(w, http.StatusNotFound, ErrCourierNotFound)
	default:
		respondInternalServerError(w, err)
	}
}

func handleAssignDeliveryError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrCouriersBusy:
		respondWithError(w, http.StatusConflict, ErrCouriersBusy)
	case usecase.ErrNoOrderID:
		respondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case usecase.ErrOrderIDExists:
		respondWithError(w, http.StatusConflict, ErrOrderIDExists)
	default:
		respondInternalServerError(w, err)
	}
}

func handleUnassignDeliveryError(w http.ResponseWriter, err error) {
	switch err {
	case usecase.ErrNoOrderID:
		respondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case usecase.ErrOrderIDNotFound:
		respondWithError(w, http.StatusNotFound, ErrOrderIDNotFound)
	default:
		respondInternalServerError(w, err)
	}
}