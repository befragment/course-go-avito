package database

import (
	"strings"
)

const (
	IDColumn            = "id"
	NameColumn          = "name"
	PhoneColumn         = "phone"
	StatusColumn        = "status"
	TransportTypeColumn = "transport_type"
	CreatedAtColumn     = "created_at"
	UpdatedAtColumn     = "updated_at"
	OrderIDColumn       = "order_id"
	AssignedAtColumn    = "assigned_at"
	DeadlineColumn      = "deadline"
	CourierIDColumn     = "courier_id"
	
	CourierTable        = "couriers"
	DeliveryTable       = "delivery"

	StatusBusy      = "busy"
	StatusAvailable = "available"

	CourierID            = CourierTable + "." + IDColumn
	CourierName          = CourierTable + "." + NameColumn
	CourierPhone         = CourierTable + "." + PhoneColumn
	CourierStatus        = CourierTable + "." + StatusColumn
	CourierTransportType = CourierTable + "." + TransportTypeColumn

	DeliveryID      = DeliveryTable + "." + IDColumn
	DeliveryOrderID = DeliveryTable + "." + OrderIDColumn
	DeliveryCourierID = DeliveryTable + "." + CourierIDColumn

	CountAll = "count(*)"
)

func BuildReturningStatement(args ...string) string {
	if len(args) == 0 {
		return ""
	}
	return "RETURNING " + strings.Join(args, ", ")
}
