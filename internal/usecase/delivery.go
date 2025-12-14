package usecase

import (
	"errors"
	"time"
	"context"
	"courier-service/internal/model"
	"courier-service/internal/repository"
)

type DelieveryUseCase struct {
	courierRepository  сourierRepository
	deliveryRepository deliveryRepository
	txRunner           txRunner
	factory deliveryCalculatorFactory
}

func NewDelieveryUseCase(
	courierRepository сourierRepository,
	deliveryRepository deliveryRepository,
	txRunner txRunner,
	factory deliveryCalculatorFactory,
) *DelieveryUseCase {
	return &DelieveryUseCase{
		courierRepository:  courierRepository,
		deliveryRepository: deliveryRepository,
		txRunner:           txRunner,
		factory: factory,
	}
}

func (u *DelieveryUseCase) AssignDelivery(ctx context.Context, req DeliveryAssignRequest) (DeliveryAssignResponse, error) {
	if req.OrderID == "" {
		return DeliveryAssignResponse{}, ErrNoOrderID
	}
	var resp DeliveryAssignResponse
	courier, delivery, err := createAssignmentInTransaction(ctx, u.txRunner, u.courierRepository, u.deliveryRepository, u.factory, req.OrderID)
	if err != nil {
		return DeliveryAssignResponse{}, err
	}
	resp = deliveryAssignResponse(courier, delivery)
	return resp, nil
}

func (u *DelieveryUseCase) UnassignDelivery(ctx context.Context, req DeliveryUnassignRequest) (DeliveryUnassignResponse, error) {
	if req.OrderID == "" {
		return DeliveryUnassignResponse{}, ErrNoOrderID
	}

	var resp DeliveryUnassignResponse
	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		couriersDelivery, err := u.deliveryRepository.CouriersDelivery(txCtx, req.OrderID)
		if err != nil {
			if errors.Is(err, repository.ErrOrderIDNotFound) {
				return ErrOrderIDNotFound
			}
			return err
		}

		if err := u.deliveryRepository.DeleteDelivery(txCtx, req.OrderID); err != nil {
			if errors.Is(err, repository.ErrOrderIDNotFound) {
				return ErrOrderIDNotFound
			}
			return err
		}

		courier, err := u.courierRepository.GetCourierById(txCtx, couriersDelivery.CourierID)
		if err != nil {
			return err
		}

		courier.ChangeStatus(model.CourierStatusAvailable)
		if err := u.courierRepository.UpdateCourier(txCtx, courier); err != nil {
			return err
		}

		resp = deliveryUnassignResponse(courier, couriersDelivery)

		return nil
	})
	if err != nil {
		return DeliveryUnassignResponse{}, err
	}

	return resp, nil
}

func createAssignmentInTransaction(
    ctx context.Context,
    txRunner txRunner,
    courierRepo сourierRepository,
    deliveryRepo deliveryRepository,
    factory deliveryCalculatorFactory,
    orderID string,
) (model.Courier, model.Delivery, error) {
    var (
        courier  model.Courier
        delivery model.Delivery
    )

    err := txRunner.Run(ctx, func(txCtx context.Context) error {
        c, err := courierRepo.FindAvailableCourier(txCtx)
        if err != nil {
            if errors.Is(err, repository.ErrCouriersBusy) {
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

        d, err := deliveryRepo.CreateDelivery(txCtx, deliveryDomain)
        if err != nil {
            if errors.Is(err, repository.ErrOrderIDExists) {
                return ErrOrderIDExists
            }
            return err
        }

        c.ChangeStatus(model.CourierStatusBusy)
        if err := courierRepo.UpdateCourier(txCtx, c); err != nil {
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