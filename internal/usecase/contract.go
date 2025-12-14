//go:generate mockgen -source=contract.go -destination=./mocks/repository_mock.go -package=mocks
package usecase

import (
	"context"
	"time"
	"courier-service/internal/model"
)

type сourierRepository interface {
	GetCourierById(ctx context.Context, id int64) (model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, courier model.Courier) (int64, error)
	UpdateCourier(ctx context.Context, courier model.Courier) error
	FindAvailableCourier(ctx context.Context) (model.Courier, error)
	ExistsCourierByPhone(ctx context.Context, phone string) (bool, error)
	FreeCouriersWithInterval(ctx context.Context) error
}

type deliveryRepository interface {
	CreateDelivery(ctx context.Context, delivery model.Delivery) (model.Delivery, error)
	CouriersDelivery(ctx context.Context, orderID string) (model.Delivery, error)
	DeleteDelivery(ctx context.Context, orderID string) error
}

type txRunner interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}

type deliveryCalculatorFactory interface {
	GetDeliveryCalculator(courierType model.CourierTransportType) DeliveryCalculator
}

type DeliveryCalculator interface {
	CalculateDeadline() time.Time
}

type orderGateway interface {
	GetOrders(ctx context.Context, from time.Time) ([]model.Order, error)
}

