package usecase

import (
	"time"
	"courier-service/internal/model"
)

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

type OnFootCalculator struct{}
type ScooterCalculator struct{}
type CarCalculator struct{}

func (c *OnFootCalculator) CalculateDeadline() time.Time {
	return time.Now().Add(15 * time.Minute)
}

func (c *ScooterCalculator) CalculateDeadline() time.Time {
	return time.Now().Add(10 * time.Minute)
}

func (c *CarCalculator) CalculateDeadline() time.Time {
	return time.Now().Add(5 * time.Minute)
}

func (f Factory) GetDeliveryCalculator(courierType model.CourierTransportType) DeliveryCalculator {
	switch courierType {
	case model.TransportTypeOnFoot:
		return &OnFootCalculator{}
	case model.TransportTypeScooter:
		return &ScooterCalculator{}
	case model.TransportTypeCar:
		return &CarCalculator{}
	default:
		return nil
	}
}