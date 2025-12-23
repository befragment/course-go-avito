package courier

import (
	"net/http"
	"courier-service/internal/usecase/courier"
	"courier-service/internal/handlers/utils"
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
)

func handleCreateError(w http.ResponseWriter, err error) {
	switch err {
	case courier.ErrInvalidCreate:
		utils.RespondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case courier.ErrInvalidPhoneNumber:
		utils.RespondWithError(w, http.StatusBadRequest, ErrInvalidPhoneNumber)
	case courier.ErrUnknownTransportType:
		utils.RespondWithError(w, http.StatusBadRequest, ErrUnknownTransportType)
	case courier.ErrPhoneNumberExists:
		utils.RespondWithError(w, http.StatusConflict, ErrPhoneAlreadyExists)
	default:
		utils.RespondInternalServerError(w, err)
	}
}

func handleUpdateError(w http.ResponseWriter, err error) {
	switch err {
	case courier.ErrInvalidUpdate:
		utils.RespondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case courier.ErrInvalidPhoneNumber:
		utils.RespondWithError(w, http.StatusBadRequest, ErrInvalidPhoneNumber)
	case courier.ErrUnknownTransportType:
		utils.RespondWithError(w, http.StatusBadRequest, ErrUnknownTransportType)
	case courier.ErrPhoneNumberExists:
		utils.RespondWithError(w, http.StatusConflict, ErrPhoneAlreadyExists)
	case courier.ErrCourierNotFound:
		utils.RespondWithError(w, http.StatusNotFound, ErrCourierNotFound)
	default:
		utils.RespondInternalServerError(w, err)
	}
}
