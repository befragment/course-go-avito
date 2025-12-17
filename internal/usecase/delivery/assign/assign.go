package assign

import (
	"context"
	"courier-service/internal/model"
	courierrepoerrors "courier-service/internal/repository/courier"
	deliveryrepoerrors "courier-service/internal/repository/delivery"
	"errors"
	"time"
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

func (u *AssignDelieveryUseCase) Assign(ctx context.Context, req DeliveryAssignRequest) (DeliveryAssignResponse, error) {
	if req.OrderID == "" {
		return DeliveryAssignResponse{}, ErrNoOrderID
	}
	var resp DeliveryAssignResponse
	courier, delivery, err := createAssignmentInTransaction(
		ctx,
		u.txRunner,
		u.courierRepository,
		u.deliveryRepository,
		u.factory,
		req.OrderID,
	)
	if err != nil {
		return DeliveryAssignResponse{}, err
	}
	resp = deliveryAssignResponse(courier, delivery)
	return resp, nil
}

func createAssignmentInTransaction(
	ctx context.Context,
	txRunner txRunner,
	courierRepository courierRepository,
	deliveryRepository deliveryRepository,
	factory deliveryCalculatorFactory,
	orderID string,
) (model.Courier, model.Delivery, error) {
	var (
		courier  model.Courier
		delivery model.Delivery
	)

	err := txRunner.Run(ctx, func(txCtx context.Context) error {
		c, err := courierRepository.FindAvailableCourier(txCtx)
		if err != nil {
			if errors.Is(err, courierrepoerrors.ErrCouriersBusy) {
				return ErrCouriersBusy
			}
			return err
		}

		dc := factory.GetDeliveryCalculator(c.TransportType)
		if dc == nil {
			return ErrUnknownTransportType
		}
		deliveryDomain := model.Delivery{
			OrderID:    orderID,
			CourierID:  c.ID,
			AssignedAt: time.Now(),
			Deadline:   dc.CalculateDeadline(),
		}

		d, err := deliveryRepository.CreateDelivery(txCtx, deliveryDomain)
		if err != nil {
			if errors.Is(err, deliveryrepoerrors.ErrOrderIDExists) {
				return ErrOrderIDExists
			}
			return err
		}

		c.ChangeStatus(model.CourierStatusBusy)
		if err := courierRepository.UpdateCourier(txCtx, c); err != nil {
			return err
		}

		courier = c
		delivery = d
		return nil
	})
	if err != nil {
		return model.Courier{}, model.Delivery{}, err
	}

	return courier, delivery, nil
}
