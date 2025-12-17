package courier

import (
	"context"
	"courier-service/internal/core"
	"courier-service/internal/model"
	courierRepo "courier-service/internal/repository/courier"
	"errors"
	"log"
	"regexp"
	"time"
)

type CourierUseCase struct {
	repository courierRepository
	factory    deliveryCalculatorFactory
}

func NewCourierUseCase(repository courierRepository, factory deliveryCalculatorFactory) *CourierUseCase {
	return &CourierUseCase{
		repository: repository,
		factory:    factory,
	}
}

func (u *CourierUseCase) CheckFreeCouriersWithInterval(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			if err := u.repository.FreeCouriersWithInterval(ctx); err != nil {
				log.Printf("Failed to check free couriers: %v", err)
			}
			log.Printf("Checked free couriers at %s", t.Format(time.RFC3339))
		}
	}
}

func (u *CourierUseCase) GetCourierById(ctx context.Context, id int64) (model.Courier, error) {
	courier, err := u.repository.GetCourierById(ctx, id)

	if err != nil {
		if errors.Is(err, courierRepo.ErrCourierNotFound) {
			return model.Courier{}, ErrCourierNotFound
		}
		return model.Courier{}, err
	}

	return courier, nil
}

func (u *CourierUseCase) GetAllCouriers(ctx context.Context) ([]model.Courier, error) {
	couriers, err := u.repository.GetAllCouriers(ctx)
	if err != nil {
		return nil, err
	}
	return couriers, nil
}

func (u *CourierUseCase) CreateCourier(ctx context.Context, courier model.Courier) (int64, error) {
	if courier.Name == "" || courier.Phone == "" || courier.Status == "" || courier.TransportType == "" {
		return 0, ErrInvalidCreate
	}

	if u.factory.GetDeliveryCalculator(courier.TransportType) == nil {
		return 0, ErrUnknownTransportType
	}

	if !ValidPhoneNumber(courier.Phone) {
		return 0, ErrInvalidPhoneNumber
	}

	if phoneExists, err := u.repository.ExistsCourierByPhone(ctx, courier.Phone); phoneExists {
		return 0, ErrPhoneNumberExists
	} else if err != nil {
		return 0, err
	}

	return u.repository.CreateCourier(ctx, courier)
}

func (u *CourierUseCase) UpdateCourier(ctx context.Context, courier model.Courier) error {
	if courier.Name == "" && courier.Phone == "" && courier.Status == "" && courier.TransportType == "" {
		return ErrInvalidUpdate
	}
	if courier.TransportType != "" {
		if u.factory.GetDeliveryCalculator(courier.TransportType) == nil {
			return ErrUnknownTransportType
		}
	}
	if courier.Phone != "" {
		if !ValidPhoneNumber(courier.Phone) {
			return ErrInvalidPhoneNumber
		}
		if phoneExists, err := u.repository.ExistsCourierByPhone(ctx, courier.Phone); err != nil {
			return err
		} else if phoneExists {
			return ErrPhoneNumberExists
		}
	}

	if err := u.repository.UpdateCourier(ctx, courier); err != nil {
		if errors.Is(err, courierRepo.ErrCourierNotFound) {
			return ErrCourierNotFound
		}
		return err
	}
	return nil
}

func ValidPhoneNumber(phone string) bool {
	return regexp.MustCompile(core.PhoneRegex).MatchString(phone)
}
