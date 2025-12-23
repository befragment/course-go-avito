package entity

import (
	"time"
	"courier-service/internal/model"
)

type CourierDB struct {
	ID        int64     	`db:"id"`
	Name      string    	`db:"name"`
	Phone     string    	`db:"phone"`
	Status    string    	`db:"status"`
	TransportType string 	`db:"transport_type"`
	CreatedAt time.Time 	`db:"created_at"`
	UpdatedAt time.Time 	`db:"updated_at"`
}

func (c CourierDB) ToModel() model.Courier {
	return model.Courier{
		ID: c.ID,
		Name: c.Name,
		Phone: c.Phone,
		Status: model.CourierStatus(c.Status),
		TransportType: model.CourierTransportType(c.TransportType),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}