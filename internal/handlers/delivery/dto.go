package delivery


import (
	assign "courier-service/internal/usecase/delivery/assign"
	"time"
)

const (
	UnassignedStatus = "unassigned"
)

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

func ToUnassignCourierResponse(courierID int64, orderID string) DeliveryUnassignResponseDTO {
	return DeliveryUnassignResponseDTO{
		OrderID:   orderID,
		Status:    UnassignedStatus,
		CourierID: courierID,
	}
}

func ToAssignCourierResponse(delivery assign.DeliveryAssignResponse) DeliveryAssignResponseDTO {
	return DeliveryAssignResponseDTO{
		CourierID:     delivery.CourierID,
		OrderID:       delivery.OrderID,
		TransportType: delivery.TransportType,
		Deadline:      delivery.Deadline,
	}
}