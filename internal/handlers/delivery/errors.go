package delivery

import (
	"net/http"
	assign "courier-service/internal/usecase/delivery/assign"
	unassign "courier-service/internal/usecase/delivery/unassign"
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
	ErrCouriersBusy          = "All couriers are busy"
	ErrOrderIDExists         = "Order id already exists"
	ErrNoCourierForOrder     = "No courier found for the order"
	ErrOrderIDNotFound       = "Order id not found"
)

func handleAssignDeliveryError(w http.ResponseWriter, err error) {
	switch err {
	case assign.ErrCouriersBusy:
		utils.RespondWithError(w, http.StatusConflict, ErrCouriersBusy)
	case assign.ErrNoOrderID:
		utils.RespondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case assign.ErrOrderIDExists:
		utils.RespondWithError(w, http.StatusConflict, ErrOrderIDExists)
	default:
		utils.RespondInternalServerError(w, err)
	}
}

func handleUnassignDeliveryError(w http.ResponseWriter, err error) {
	switch err {
	case unassign.ErrNoOrderID:
		utils.RespondWithError(w, http.StatusBadRequest, ErrMissingRequiredFields)
	case unassign.ErrOrderIDNotFound:
		utils.RespondWithError(w, http.StatusNotFound, ErrOrderIDNotFound)
	default:
		utils.RespondInternalServerError(w, err)
	}
}