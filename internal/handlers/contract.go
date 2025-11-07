package handlers

import (
	"context"
	"courier-service/internal/model"
)

type CourierUseCase interface {
	GetById(ctx context.Context, id int64) (*model.Courier, error)
	GetAll(ctx context.Context) ([]model.Courier, error)
	Create(ctx context.Context, req *model.CourierCreateRequest) (int64, error)
	Update(ctx context.Context, req *model.CourierUpdateRequest) error
}