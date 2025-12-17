//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package courier

import (
	"context"
	"courier-service/internal/model"
	utils "courier-service/internal/usecase/utils"
)

type courierRepository interface {
	GetCourierById(ctx context.Context, id int64) (model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, courier model.Courier) (int64, error)
	UpdateCourier(ctx context.Context, courier model.Courier) error
	FindAvailableCourier(ctx context.Context) (model.Courier, error)
	ExistsCourierByPhone(ctx context.Context, phone string) (bool, error)
	FreeCouriersWithInterval(ctx context.Context) error
}

// DeliveryCalculator aliases the shared utils.DeliveryCalculator
// so there is a single interface type used across usecases and mocks.
type DeliveryCalculator = utils.DeliveryCalculator

type deliveryCalculatorFactory interface {
	GetDeliveryCalculator(courierType model.CourierTransportType) DeliveryCalculator
}
