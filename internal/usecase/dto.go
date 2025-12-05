package usecase

import (
	"courier-service/internal/model"
	"courier-service/internal/repository"
	"time"
)

func courierUpdateRequestToCourierDB(req model.CourierUpdateRequest) repository.CourierDB {
	out := repository.CourierDB{ID: req.ID}
	if req.Name != nil {
		out.Name = *req.Name
	}
	if req.Phone != nil {
		out.Phone = *req.Phone
	}
	if req.Status != nil {
		out.Status = *req.Status
	}
	if req.TransportType != nil {
		out.TransportType = *req.TransportType
	}
	return out
}

func courierCreateRequestToCourierDB(req model.CourierCreateRequest) repository.CourierDB {
	return repository.CourierDB{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}
}

func deliveryAssignRequestToDeliveryDB(orderID string, courierDB repository.CourierDB) (model.DeliveryDB, error) {
	duration, err := transportTypeTime(courierDB.TransportType)
	if err != nil {
		return model.DeliveryDB{}, err
	}
	return model.DeliveryDB{
		OrderID:    orderID,
		CourierID:  courierDB.ID,
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(duration),
	}, nil
}

func delieveryAssignResponse(courier model.Courier, delivery model.Delivery) model.DeliveryAssignResponse {
	return model.DeliveryAssignResponse{
		CourierID:     courier.ID,
		OrderID:       delivery.OrderID,
		TransportType: courier.TransportType,
		Deadline:      delivery.Deadline,
	}
}

func delieveryUnassignResponse(courier model.Courier, delivery model.Delivery) model.DeliveryUnassignResponse {
	return model.DeliveryUnassignResponse{
		OrderID:   delivery.OrderID,
		Status:    "unassigned",
		CourierID: courier.ID,
	}
}
