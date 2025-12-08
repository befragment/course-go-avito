package repository

import "time"

type CourierDB struct {
	ID        int64     	`json:"id"`
	Name      string    	`json:"name"`
	Phone     string    	`json:"phone"`
	Status    string    	`json:"status"`
	TransportType string 	`json:"transport_type"`
	CreatedAt time.Time 	`json:"created_at,omitempty"`
	UpdatedAt time.Time 	`json:"updated_at,omitempty"`
}

type DeliveryDB struct {
	ID         int64     `json:"id"`
	CourierID  int64     `json:"courier_id"`
	OrderID    string    `json:"order_id,omitempty"`
	AssignedAt time.Time `json:"assigned_at,omitempty"`
	Deadline   time.Time `json:"deadline,omitempty"`
}