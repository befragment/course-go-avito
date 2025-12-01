package handlers

import (
	"context"
	"courier-service/internal/model"
)

type сourierUseCase interface {
	GetCourierById(ctx context.Context, id int64) (*model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, req *model.CourierCreateRequest) (int64, error)
	UpdateCourier(ctx context.Context, req *model.CourierUpdateRequest) error
}

type deliveryUseCase interface {
	AssignDelivery(ctx context.Context, req *model.DeliveryAssignRequest) (model.DeliveryAssignResponse, error)
	UnassignDelivery(ctx context.Context, req *model.DeliveryUnassignRequest) (model.DeliveryUnassignResponse, error)
}
