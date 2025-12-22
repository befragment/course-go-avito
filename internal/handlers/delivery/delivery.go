package delivery

import (
	"courier-service/internal/handlers/utils"
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

	assignment, err := c.assign.Assign(ctx, req.OrderID)
	response := ToAssignCourierResponse(assignment)

	if err != nil {
		handleAssignDeliveryError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (c *DeliveryController) UnassignDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req DeliveryUnassignRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	courierID, err := c.unassign.Unassign(ctx, req.OrderID)
	response := ToUnassignCourierResponse(courierID, req.OrderID)

	if err != nil {
		handleUnassignDeliveryError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
