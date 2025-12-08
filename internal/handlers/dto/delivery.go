package dto

import "time"

type DeliveryAssignRequestDTO struct {
    OrderID string `json:"order_id"`
}

type DeliveryUnassignRequestDTO struct {
    OrderID string `json:"order_id"`
}

type DeliveryAssignResponseDTO struct {
	CourierID     int64     `json:"courier_id"`
	OrderID       string    `json:"order_id"`
	TransportType string    `json:"transport_type"`
	Deadline      time.Time `json:"delivery_deadline"`
}

type DeliveryUnassignResponseDTO struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CourierID int64  `json:"courier_id"`
}

