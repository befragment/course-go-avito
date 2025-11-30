package model

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

type Courier struct {
	ID        int64     	`json:"id"`
	Name      string    	`json:"name"`
	Phone     string    	`json:"phone"`
	Status    string    	`json:"status"`
	TransportType string 	`json:"transport_type"`
	CreatedAt time.Time 	`json:"-"`
	UpdatedAt time.Time 	`json:"-"`
}

type CourierCreateRequest struct {
	Name   			string 	`json:"name"`
	TransportType 	string 	`json:"transport_type"`
	Phone  			string 	`json:"phone"`
	Status 			string 	`json:"status"`
}

type CourierUpdateRequest struct {
	ID     			int64  	`json:"id"`
	TransportType 	*string `json:"transport_type"`
	Name   			*string `json:"name"`
	Phone  			*string `json:"phone"`
	Status 			*string `json:"status"`
}
