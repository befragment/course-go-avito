//go:generate mockgen -source=contract.go -destination=./mocks/repository_mock.go -package=mocks
package usecase

import (
	"context"
	"courier-service/internal/model"
)

type сourierRepository interface {
	GetCourierById(ctx context.Context, id int64) (*model.CourierDB, error)
	GetAllCouriers(ctx context.Context) ([]model.CourierDB, error)
	CreateCourier(ctx context.Context, courier *model.CourierDB) (int64, error)
	UpdateCourier(ctx context.Context, courier *model.CourierDB) error
	FindAvailableCourier(ctx context.Context) (*model.CourierDB, error)
	ExistsCourierByPhone(ctx context.Context, phone string) (bool, error)
	FreeCouriersWithInterval(ctx context.Context) error
}

type deliveryRepository interface {
	CreateDelivery(ctx context.Context, delivery *model.DeliveryDB) (*model.Delivery, error)
	CouriersDelivery(ctx context.Context, orderID string) (*model.DeliveryDB, error)
	DeleteDelivery(ctx context.Context, orderID string) error
}

type txRunner interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
