package usecase

import (
	"time"
	"courier-service/internal/model"
)

type DeliveryAssignRequest struct {
	OrderID string `json:"order_id"`
}

type DeliveryUnassignRequest struct {
	OrderID string `json:"order_id"`
}

type DeliveryAssignResponse struct {
	CourierID     int64     `json:"courier_id"`
	OrderID       string    `json:"order_id"`
	TransportType string    `json:"transport_type"`
	Deadline      time.Time `json:"delivery_deadline"`
}

type DeliveryUnassignResponse struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CourierID int64  `json:"courier_id"`
}

func deliveryAssignResponse(courier model.Courier, delivery model.Delivery) DeliveryAssignResponse {
	return DeliveryAssignResponse{
		CourierID:     courier.ID,
		OrderID:       delivery.OrderID,
		TransportType: courier.TransportType,
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
