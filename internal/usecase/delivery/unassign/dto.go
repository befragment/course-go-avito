package unassign

import (
	"courier-service/internal/model"
)

type DeliveryUnassignRequest struct {
	OrderID string
}

type DeliveryUnassignResponse struct {
	OrderID   string 
	Status    string 
	CourierID int64  
}


func deliveryUnassignResponse(courier model.Courier, delivery model.Delivery) DeliveryUnassignResponse {
	return DeliveryUnassignResponse{
		OrderID:   delivery.OrderID,
		Status:    "unassigned",
		CourierID: courier.ID,
	}
}
