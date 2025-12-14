package usecase

import (
	"time"
	"courier-service/internal/model"
)

type DeliveryAssignRequest struct {
	OrderID string
}

type DeliveryUnassignRequest struct {
	OrderID string
}

type DeliveryAssignResponse struct {
	CourierID     int64     
	OrderID       string    
	TransportType string    
	Deadline      time.Time 
}

type DeliveryUnassignResponse struct {
	OrderID   string 
	Status    string 
	CourierID int64  
}

func deliveryAssignResponse(courier model.Courier, delivery model.Delivery) DeliveryAssignResponse {
	return DeliveryAssignResponse{
		CourierID:     courier.ID,
		OrderID:       delivery.OrderID,
		TransportType: string(courier.TransportType),
		Deadline:      delivery.Deadline,
	}
}

func deliveryUnassignResponse(courier model.Courier, delivery model.Delivery) DeliveryUnassignResponse {
	return DeliveryUnassignResponse{
		OrderID:   delivery.OrderID,
		Status:    "unassigned",
		CourierID: courier.ID,
	}
}
