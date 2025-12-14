package model

import "time"

type Order struct {
    ID           string
	UserID       string
    OrderNumber  string
	FIO          string
    RestaurantID string
    Items        []Item
    TotalPrice   int64
    Address      Address
    Status       OrderStatus
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type OrderStatus string

const (
    OrderStatusCreated      OrderStatus = "created"
    OrderStatusPending      OrderStatus = "pending"
    OrderStatusConfirmed    OrderStatus = "confirmed"
    OrderStatusCooking      OrderStatus = "cooking"
    OrderStatusDelivering   OrderStatus = "delivering"
    OrderStatusCompleted    OrderStatus = "completed"
	OrderStatusDeleted      OrderStatus = "deleted"
	OrderStatusCancelled    OrderStatus = "cancelled"
	OrderStatusFailed       OrderStatus = "failed"
)

type Item struct {
	FoodID   string
    Name     string
    Price    int64
    Quantity int64
}

type Address struct {
    Street    string
    House     string
    Apartment string
    Floor     string
    Comment   string
}