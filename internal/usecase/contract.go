package usecase

import (
	"context"
	"courier-service/internal/model"
)

type сourierRepository interface {
	GetById(ctx context.Context, id int64) (*model.CourierDB, error)
	GetAll(ctx context.Context) ([]model.CourierDB, error)
	Create(ctx context.Context, courier *model.CourierDB) (int64, error)
	Update(ctx context.Context, courier *model.CourierDB) error
	FindAvailiable(ctx context.Context) (*model.CourierDB, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	FreeCouriers(ctx context.Context) error
}

type deliveryRepository interface {
	Create(ctx context.Context, delivery *model.DeliveryDB) (*model.Delivery, error)
	CouriersDelivery(ctx context.Context, orderID string) (*model.DeliveryDB, error)
	Delete(ctx context.Context, orderID string) error
}

type txRunner interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
