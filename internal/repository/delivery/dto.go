package delivery

import (
	"time"
)

type DeliveryDB struct {
	ID         int64     `db:"id"`
	CourierID  int64     `db:"courier_id"`
	OrderID    string    `db:"order_id,omitempty"`
	AssignedAt time.Time `db:"assigned_at,omitempty"`
	Deadline   time.Time `db:"deadline,omitempty"`
}