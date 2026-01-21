package assign

import (
	"context"
	"errors"
	"time"

	"courier-service/internal/model"
	courierrepoerrors "courier-service/internal/repository/courier"
	deliveryrepoerrors "courier-service/internal/repository/delivery"
)

type AssignDelieveryUseCase struct {
	courierRepository  courierRepository
	deliveryRepository deliveryRepository
	txRunner           txRunner
	factory            deliveryCalculatorFactory
}

func NewAssignDelieveryUseCase(
	courierRepository courierRepository,
	deliveryRepository deliveryRepository,
	txRunner txRunner,
	factory deliveryCalculatorFactory,
) *AssignDelieveryUseCase {
	return &AssignDelieveryUseCase{
		courierRepository:  courierRepository,
		deliveryRepository: deliveryRepository,
		txRunner:           txRunner,
		factory:            factory,
	}
}

func (u *AssignDelieveryUseCase) Assign(ctx context.Context, OrderID string) (DeliveryAssignResponse, error) {
	if OrderID == "" {
		return DeliveryAssignResponse{}, ErrNoOrderID
	}
	var resp DeliveryAssignResponse
	var courier model.Courier
	var delivery model.Delivery
	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		c, err := u.courierRepository.FindAvailableCourier(txCtx)
		if err != nil {
			if errors.Is(err, courierrepoerrors.ErrCouriersBusy) {
				return ErrCouriersBusy
			}
			return err
		}

		dc := u.factory.GetDeliveryCalculator(c.TransportType)
		if dc == nil {
			return ErrUnknownTransportType
		}
		deliveryDomain := model.Delivery{
			OrderID:    OrderID,
			CourierID:  c.ID,
			AssignedAt: time.Now(),
			Deadline:   dc.CalculateDeadline(),
		}

		d, err := u.deliveryRepository.CreateDelivery(txCtx, deliveryDomain)
		if err != nil {
			if errors.Is(err, deliveryrepoerrors.ErrOrderIDExists) {
				return ErrOrderIDExists
			}
			return err
		}

		c.ChangeStatus(model.CourierStatusBusy)
		if err := u.courierRepository.UpdateCourier(txCtx, c); err != nil {
			return err
		}

		courier = c
		delivery = d
		return nil
	})
	if err != nil {
		return DeliveryAssignResponse{}, err
	}
	resp = deliveryAssignResponse(courier, delivery)
	return resp, nil
}
