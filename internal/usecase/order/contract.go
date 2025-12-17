//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package order

import (
	"context"
	"courier-service/internal/model"
	utils "courier-service/internal/usecase/utils"
	"time"
)

// DeliveryCalculator aliases the shared utils.DeliveryCalculator
// so that all usecases share the same delivery-time interface.
type DeliveryCalculator = utils.DeliveryCalculator

type deliveryCalculatorFactory interface {
	GetDeliveryCalculator(courierType model.CourierTransportType) DeliveryCalculator
}

type orderGateway interface {
	GetOrders(ctx context.Context, from time.Time) ([]model.Order, error)
}

type courierRepository interface {
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
