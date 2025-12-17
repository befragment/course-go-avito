package delivery

import (
	"courier-service/internal/handlers/utils"
	assign "courier-service/internal/usecase/delivery/assign"
	unassign "courier-service/internal/usecase/delivery/unassign"
	"encoding/json"
	"net/http"
)

type DeliveryController struct {
	assign   assignUsecase
	unassign unassignUsecase
}

func NewDeliveryController(assign assignUsecase, unassign unassignUsecase) *DeliveryController {
	return &DeliveryController{assign: assign, unassign: unassign}
}

func (c *DeliveryController) AssignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req DeliveryAssignRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	delivery, err := c.assign.Assign(ctx, assign.DeliveryAssignRequest{
		OrderID: req.OrderID,
	})

	if err != nil {
		handleAssignDeliveryError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, delivery)
}

func (c *DeliveryController) UnassignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req DeliveryUnassignRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	delivery, err := c.unassign.Unassign(ctx, unassign.DeliveryUnassignRequest{
		OrderID: req.OrderID,
	})

	if err != nil {
		handleUnassignDeliveryError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, delivery)
}
