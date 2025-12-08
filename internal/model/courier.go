package model

import "time"

type Courier struct {
	ID        int64     	`json:"id"`
	Name      string    	`json:"name"`
	Phone     string    	`json:"phone"`
	Status    string    	`json:"status"`
	TransportType string 	`json:"transport_type"`
	CreatedAt time.Time 	`json:"-"`
	UpdatedAt time.Time 	`json:"-"`
}