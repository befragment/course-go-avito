package handlers

import (
	"net/http"
	"encoding/json"
	"courier-service/internal/model"
)

type DeliveryController struct {
	usecase deliveryUseCase
}

func NewDeliveryController(usecase deliveryUseCase) *DeliveryController {
	return &DeliveryController{usecase: usecase}
}

func (c *DeliveryController) AssignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.DeliveryAssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return 
	}
	
	delivery, err := c.usecase.AssignDelivery(ctx, &req)
	if err != nil {
		handleAssignDeliveryError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, delivery)
}

func (c *DeliveryController) UnassignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.DeliveryUnassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	delivery, err := c.usecase.UnassignDelivery(ctx, &req)
	if err != nil {
		handleUnassignDeliveryError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, delivery)
}