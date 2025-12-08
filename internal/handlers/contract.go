package handlers

import (
	"context"
	"courier-service/internal/model"
	"courier-service/internal/usecase"
)

type сourierUseCase interface {
	GetCourierById(ctx context.Context, id int64) (model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, courier model.Courier) (int64, error)
	UpdateCourier(ctx context.Context, courier model.Courier) error
}

type deliveryUseCase interface {
	AssignDelivery(ctx context.Context, req usecase.DeliveryAssignRequest) (usecase.DeliveryAssignResponse, error)
	UnassignDelivery(ctx context.Context, req usecase.DeliveryUnassignRequest) (usecase.DeliveryUnassignResponse, error)
}
