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
}

func NewDelieveryUseCase(
	courierRepository сourierRepository,
	deliveryRepository deliveryRepository,
	txRunner txRunner,
) *DelieveryUseCase {
	return &DelieveryUseCase{
		courierRepository:  courierRepository,
		deliveryRepository: deliveryRepository,
		txRunner:           txRunner,
	}
}

func transportTypeTime(ttype string) (time.Duration, error) {
	switch ttype {
	case "car":
		return 5 * time.Minute, nil
	case "scooter":
		return 15 * time.Minute, nil
	case "on_foot":
		return 30 * time.Minute, nil
	default:
		return 0, ErrUnknownTransportType
	}
}

func (u *DelieveryUseCase) AssignDelivery(ctx context.Context, req DeliveryAssignRequest) (DeliveryAssignResponse, error) {
	if req.OrderID == "" {
		return DeliveryAssignResponse{}, ErrNoOrderID
	}

	var resp DeliveryAssignResponse
	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		courier, err := u.courierRepository.FindAvailableCourier(txCtx)
		if err != nil {
			if errors.Is(err, repository.ErrCouriersBusy) {
				return ErrCouriersBusy
			}
			return err
		}

		duration, err := transportTypeTime(courier.TransportType)
		if err != nil {
			return err
		}

		deliveryDomain := model.Delivery{
			OrderID:    req.OrderID,
			CourierID:  courier.ID,
			AssignedAt: time.Now(),
			Deadline:   time.Now().Add(duration),
		}

		delivery, err := u.deliveryRepository.CreateDelivery(txCtx, deliveryDomain)
		if err != nil {
			if errors.Is(err, repository.ErrOrderIDExists) {
				return ErrOrderIDExists
			}
			return err
		}

		courier.Status = "busy"
		if err := u.courierRepository.UpdateCourier(txCtx, courier); err != nil {
			return err
		}
		resp = deliveryAssignResponse(courier, delivery)
		return nil
	})

	if err != nil {
		return DeliveryAssignResponse{}, err
	}
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

		courier.Status = "available"
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
