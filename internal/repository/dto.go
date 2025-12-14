package repository

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

type DeliveryDB struct {
	ID         int64     `db:"id"`
	CourierID  int64     `db:"courier_id"`
	OrderID    string    `db:"order_id,omitempty"`
	AssignedAt time.Time `db:"assigned_at,omitempty"`
	Deadline   time.Time `db:"deadline,omitempty"`
}

func courierDBToModel(c CourierDB) model.Courier {
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