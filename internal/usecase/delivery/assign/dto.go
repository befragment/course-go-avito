package assign

import (
	"time"
	"courier-service/internal/model"
)

type DeliveryAssignRequest struct {
	OrderID string
}

type DeliveryAssignResponse struct {
	CourierID     int64     
	OrderID       string    
	TransportType string    
	Deadline      time.Time 
}

func deliveryAssignResponse(courier model.Courier, delivery model.Delivery) DeliveryAssignResponse {
	return DeliveryAssignResponse{
		CourierID:     courier.ID,
		OrderID:       delivery.OrderID,
		TransportType: string(courier.TransportType),
		Deadline:      delivery.Deadline,
	}
}