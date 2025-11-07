package usecase

import (
	"context"
	"courier-service/internal/model"
)

type CourierRepository interface {
	GetById(ctx context.Context, id int64) (*model.CourierDB, error)
	GetAll(ctx context.Context) ([]model.CourierDB, error)
	Create(ctx context.Context, courier *model.CourierDB) (int64, error)
	Update(ctx context.Context, courier *model.CourierDB) error
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}
