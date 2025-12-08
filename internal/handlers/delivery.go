package handlers

import (
	"net/http"
	"encoding/json"
	"courier-service/internal/handlers/dto"
	"courier-service/internal/usecase"
)

type DeliveryController struct {
	usecase deliveryUseCase
}

func NewDeliveryController(usecase deliveryUseCase) *DeliveryController {
	return &DeliveryController{usecase: usecase}
}

func (c *DeliveryController) AssignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.DeliveryAssignRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return 
	}

	delivery, err := c.usecase.AssignDelivery(ctx, usecase.DeliveryAssignRequest{
		OrderID: req.OrderID,
	})

	if err != nil {
		handleAssignDeliveryError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, delivery)
}

func (c *DeliveryController) UnassignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.DeliveryUnassignRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	delivery, err := c.usecase.UnassignDelivery(ctx, usecase.DeliveryUnassignRequest{
		OrderID: req.OrderID,
	})
	
	if err != nil {
		handleUnassignDeliveryError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, delivery)
}
