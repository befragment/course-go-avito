package order

import (
	"time"
	assign "courier-service/internal/usecase/delivery/assign"
	complete "courier-service/internal/usecase/delivery/complete"
	unassign "courier-service/internal/usecase/delivery/unassign"
)

type orderChangedDto struct {
	OrderID string `json:"order_id"`
	Status orderStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// type orderChangedDto struct {
//     ID           string `json:"id"`
// 	UserID       string `json:"user_id"`
//     OrderNumber  string `json:"order_number"`
// 	FIO          string `json:"fio"`
//     RestaurantID string `json:"restaurant_id"`
//     Items        []Item `json:"items"`
//     TotalPrice   int64 `json:"total_price"`
//     Address      Address `json:"address"`
//     Status       orderStatus `json:"status"`
//     CreatedAt    time.Time `json:"created_at"`
//     UpdatedAt    time.Time `json:"updated_at"`
// }

type orderStatus string

const (
    OrderStatusCreated      orderStatus = "created"
    OrderStatusPending      orderStatus = "pending"
    OrderStatusConfirmed    orderStatus = "confirmed"
    OrderStatusCooking      orderStatus = "cooking"
    OrderStatusDelivering   orderStatus = "delivering"
    OrderStatusCompleted    orderStatus = "completed"
	OrderStatusDeleted      orderStatus = "deleted"
	OrderStatusCancelled    orderStatus = "cancelled"
	OrderStatusFailed       orderStatus = "failed"
)

// type Item struct {
// 	FoodID   string `json:"food_id"`
//     Name     string `json:"name"`
//     Price    int64 `json:"price"`
//     Quantity int64 `json:"quantity"`
// }

// type Address struct {
//     Street    string `json:"street"`
//     House     string `json:"house"`
//     Apartment string `json:"apartment"`
//     Floor     string `json:"floor"`
//     Comment   string `json:"comment"`
// }

func orderIDtoDeliveryAssignRequest(orderID string) assign.DeliveryAssignRequest {
	return assign.DeliveryAssignRequest{
		OrderID: orderID,
	}
}

func orderIDtoDeliveryUnassignRequest(orderID string) unassign.DeliveryUnassignRequest {
	return unassign.DeliveryUnassignRequest{
		OrderID: orderID,
	}
}

func orderIDtoCompleteDeliveryRequest(orderID string) complete.CompleteDeliveryRequest {
	return complete.CompleteDeliveryRequest{
		OrderID: orderID,
	}
}