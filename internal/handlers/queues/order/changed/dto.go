package order

import (
	"time"
)

type orderChangedDto struct {
	OrderID   string    `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
