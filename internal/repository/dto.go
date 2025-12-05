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