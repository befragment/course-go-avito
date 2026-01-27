package model

import "time"

type Order struct {
	ID        string
	CreatedAt time.Time
	Status    OrderStatus
}

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "canceled"
)
