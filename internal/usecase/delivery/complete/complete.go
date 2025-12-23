package complete

import (
	"context"
	"courier-service/internal/model"
	courierrepo "courier-service/internal/repository/courier"
	"errors"
	"fmt"
)

type CompleteDeliveryUseCase struct {
	courierRepository courierRepository
}

func NewCompleteDeliveryUseCase(courierRepository courierRepository) *CompleteDeliveryUseCase {
	return &CompleteDeliveryUseCase{courierRepository: courierRepository}
}

func (u *CompleteDeliveryUseCase) Complete(ctx context.Context, OrderID string) error {
	courierID, err := u.courierRepository.GetCourierIDByOrderID(ctx, OrderID)
	if err != nil {
		if errors.Is(err, courierrepo.ErrOrderNotFound) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("database error: %w", err)
	}

	err = u.courierRepository.UpdateCourier(ctx, model.Courier{
		ID: courierID, Status: model.CourierStatusAvailable,
	})
	if err != nil {
		return err
	}

	return nil
}
