package unassign

import (
	"context"
	"courier-service/internal/model"
	deliveryRepo "courier-service/internal/repository/delivery"
	"errors"
)

type UnassignDelieveryUseCase struct {
	courierRepository  courierRepository
	deliveryRepository deliveryRepository
	txRunner           txRunner
}

func NewUnassignDelieveryUseCase(
	courierRepository courierRepository,
	deliveryRepository deliveryRepository,
	txRunner txRunner,
) *UnassignDelieveryUseCase {
	return &UnassignDelieveryUseCase{
		courierRepository:  courierRepository,
		deliveryRepository: deliveryRepository,
		txRunner:           txRunner,
	}
}

func (u *UnassignDelieveryUseCase) Unassign(ctx context.Context, req DeliveryUnassignRequest) (DeliveryUnassignResponse, error) {
	if req.OrderID == "" {
		return DeliveryUnassignResponse{}, ErrNoOrderID
	}

	var resp DeliveryUnassignResponse
	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		couriersDelivery, err := u.deliveryRepository.CouriersDelivery(txCtx, req.OrderID)
		if err != nil {
			if errors.Is(err, deliveryRepo.ErrOrderIDNotFound) {
				return ErrOrderIDNotFound
			}
			return err
		}

		if err := u.deliveryRepository.DeleteDelivery(txCtx, req.OrderID); err != nil {
			if errors.Is(err, deliveryRepo.ErrOrderIDNotFound) {
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
