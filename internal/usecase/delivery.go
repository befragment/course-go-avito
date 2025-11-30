package usecase

import (
	"context"
	"courier-service/internal/model"
	"courier-service/internal/repository"
	"errors"
	"time"
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
		return 20 * time.Second, nil
	case "scooter":
		return 40 * time.Second, nil
	case "on_foot":
		return 60 * time.Second, nil
	default:
		return 0, ErrUnknownTransportType
	}
}

func (u *DelieveryUseCase) AssignDelivery(
	ctx context.Context,
	req *model.DeliveryAssignRequest,
) (model.DeliveryAssignResponse, error) {
	if req.OrderID == "" {
		return model.DeliveryAssignResponse{}, ErrNoOrderID
	}

	var resp model.DeliveryAssignResponse

	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		courierDB, err := u.courierRepository.FindAvailiable(txCtx)
		if err != nil {
			if errors.Is(err, repository.ErrCouriersBusy) {
				return ErrCouriersBusy
			}
			return err
		}

		deliveryDB, err := deliveryAssignRequestToDeliveryDB(req.OrderID, *courierDB)
		if err != nil {
			return err
		}

		delivery, err := u.deliveryRepository.Create(txCtx, &deliveryDB)
		if err != nil {
			if errors.Is(err, repository.ErrOrderIDExists) {
				return ErrOrderIDExists
			}
			return err
		}

		courierDB.Status = "busy"
		if err := u.courierRepository.Update(txCtx, courierDB); err != nil {
			return err
		}

		courier := model.Courier(*courierDB)
		resp = delieveryAssignResponse(courier, *delivery)
		return nil
	})

	if err != nil {
		return model.DeliveryAssignResponse{}, err
	}
	return resp, nil
}

func (u *DelieveryUseCase) UnassignDelivery(
	ctx context.Context,
	req *model.DeliveryUnassignRequest,
) (model.DeliveryUnassignResponse, error) {
	if req.OrderID == "" {
		return model.DeliveryUnassignResponse{}, ErrNoOrderID
	}

	var resp model.DeliveryUnassignResponse

	err := u.txRunner.Run(ctx, func(txCtx context.Context) error {
		couriersDelivery, err := u.deliveryRepository.CouriersDelivery(txCtx, req.OrderID)
		if err != nil {
			if errors.Is(err, repository.ErrOrderIDNotFound) {
				return ErrOrderIDNotFound
			}
			return err
		}

		if err := u.deliveryRepository.Delete(txCtx, req.OrderID); err != nil {
			if errors.Is(err, repository.ErrOrderIDNotFound) {
				return ErrOrderIDNotFound
			}
			return err
		}

		courierDB, err := u.courierRepository.GetById(txCtx, couriersDelivery.CourierID)
		if err != nil {
			return err
		}

		courierDB.Status = "available"
		if err := u.courierRepository.Update(txCtx, courierDB); err != nil {
			return err
		}

		courier := model.Courier(*courierDB)
		delivery := model.Delivery(*couriersDelivery)
		resp = delieveryUnassignResponse(courier, delivery)

		return nil
	})
	if err != nil {
		return model.DeliveryUnassignResponse{}, err
	}

	return resp, nil
}
