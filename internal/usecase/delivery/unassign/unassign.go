package unassign

import (
	"context"
	"errors"

	"courier-service/internal/model"
	deliveryRepo "courier-service/internal/repository/delivery"
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

func (u *UnassignDelieveryUseCase) Unassign(ctx context.Context, OrderID string) (int64, error) {
	if OrderID == "" {
		return 0, ErrNoOrderID
	}

	var courierID int64
	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		couriersDelivery, err := u.deliveryRepository.CouriersDelivery(txCtx, OrderID)
		if err != nil {
			if errors.Is(err, deliveryRepo.ErrOrderIDNotFound) {
				return ErrOrderIDNotFound
			}
			return err
		}

		if err := u.deliveryRepository.DeleteDelivery(txCtx, OrderID); err != nil {
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

		courierID = courier.ID

		return nil
	})
	if err != nil {
		return 0, err
	}

	return courierID, nil
}
