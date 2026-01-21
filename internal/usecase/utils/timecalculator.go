package utils

import (
	"time"

	"courier-service/internal/model"
)

type DeliveryCalculator interface {
	CalculateDeadline() time.Time
}

type TimeCalculatorFactory struct{}

func NewTimeCalculatorFactory() *TimeCalculatorFactory {
	return &TimeCalculatorFactory{}
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

func (f TimeCalculatorFactory) GetDeliveryCalculator(courierType model.CourierTransportType) DeliveryCalculator {
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
