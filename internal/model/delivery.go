package model

import "time"

type DeliveryDB struct {
	ID int64 `json:"id"`
	CourierID int64 `json:"courier_id"`
	OrderID string `json:"order_id,omitempty"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`
	Deadline time.Time `json:"deadline,omitempty"`
}

type Delivery struct {
	ID int64 `json:"id"`
	CourierID int64 `json:"courier_id"`
	OrderID string `json:"order_id,omitempty"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`
	Deadline time.Time `json:"deadline,omitempty"`
}

type DeliveryAssignRequest struct {
	OrderID string `json:"order_id"`
}

type DeliveryUnassignRequest struct {
	OrderID string `json:"order_id"`
}

type DeliveryAssignResponse struct {
	CourierID int64 `json:"courier_id"`
	OrderID string `json:"order_id"`
	TransportType string `json:"transport_type"`
	Deadline time.Time `json:"delivery_deadline"`
}

type DeliveryUnassignResponse struct {
	OrderID string `json:"order_id"`
	Status string `json:"status"`
	CourierID int64 `json:"courier_id"`
}
