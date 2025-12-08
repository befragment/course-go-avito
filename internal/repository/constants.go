package repository

import (
	"strings"
)

const (
	idColumn = "id"
	nameColumn = "name"
	phoneColumn = "phone"
	statusColumn = "status"
	transportTypeColumn = "transport_type"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
	orderIdColumn = "order_id"
	assignedAtColumn = "assigned_at"
	deadlineColumn = "deadline"
	courierIdColumn = "courier_id"
	courierTable = "couriers"
	deliveryTable = "delivery"

	statusBusy = "busy"
	statusAvailable = "available"

	courierID = courierTable + "." + idColumn
	courierName = courierTable + "." + nameColumn
	courierPhone = courierTable + "." + phoneColumn
	courierStatus = courierTable + "." + statusColumn
	courierTransportType = courierTable + "." + transportTypeColumn
	// courierCreatedAt = courierTable + "." + createdAtColumn
	// courierUpdatedAt = courierTable + "." + updatedAtColumn

	deliveryID = deliveryTable + "." + idColumn
	deliveryOrderID = deliveryTable + "." + orderIdColumn
	// deliveryAssignedAt = deliveryTable + "." + assignedAtColumn
	// deliveryDeadline = deliveryTable + "." + deadlineColumn
	deliveryCourierID = deliveryTable + "." + courierIdColumn

	countAll = "count(*)"
)

func buildReturningStatement(args ...string) string {
	if len(args) == 0 {
		return ""
	}
	return "RETURNING " + strings.Join(args, ", ")
}