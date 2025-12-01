package usecase

import (
	"context"
	"courier-service/internal/core"
	"courier-service/internal/model"
	"courier-service/internal/repository"
	"errors"
	"log"
	"regexp"
	"time"
)

type CourierUseCase struct {
	repository сourierRepository
}

func NewCourierUseCase(repository сourierRepository) *CourierUseCase {
	return &CourierUseCase{repository: repository}
}

func CheckFreeCouriers(ctx context.Context, u *CourierUseCase) {
	CheckFreeCouriersWithInterval(ctx, u, core.CheckFreeCouriersInterval)
}

func CheckFreeCouriersWithInterval(ctx context.Context, u *CourierUseCase, interval time.Duration) {
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

func (u CourierUseCase) GetCourierById(ctx context.Context, id int64) (*model.Courier, error) {
	courierDB, err := u.repository.GetCourierById(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrCourierNotFound) {
			return nil, ErrCourierNotFound
		}
		return nil, err
	}

	courier := model.Courier(*courierDB)
	return &courier, nil
}

func (u CourierUseCase) GetAllCouriers(ctx context.Context) ([]model.Courier, error) {
	couriersDB, err := u.repository.GetAllCouriers(ctx)

	if err != nil {
		return nil, err
	}

	couriers := make([]model.Courier, len(couriersDB))
	for i, c := range couriersDB {
		couriers[i] = model.Courier(c)
	}

	return couriers, nil
}

func (u CourierUseCase) CreateCourier(ctx context.Context, req *model.CourierCreateRequest) (int64, error) {
	if req.Name == "" || req.Phone == "" || req.Status == "" || req.TransportType == "" {
		return 0, ErrInvalidCreate
	}

	if _, err := transportTypeTime(req.TransportType); err != nil {
		return 0, ErrUnknownTransportType
	}

	if !ValidPhoneNumber(req.Phone) {
		return 0, ErrInvalidPhoneNumber
	}

	if phoneExists, err := u.repository.ExistsCourierByPhone(ctx, req.Phone); phoneExists {
		return 0, ErrPhoneNumberExists
	} else if err != nil {
		return 0, err
	}

	courierDB := courierCreateRequestToCourierDB(*req)
	return u.repository.CreateCourier(ctx, &courierDB)
}

func (u CourierUseCase) UpdateCourier(ctx context.Context, req *model.CourierUpdateRequest) error {
	if req.Name == nil && req.Phone == nil && req.Status == nil && req.TransportType == nil {
		return ErrInvalidUpdate
	}
	if req.TransportType != nil {
		if _, err := transportTypeTime(*req.TransportType); err != nil {
			return ErrUnknownTransportType
		}
	}
	if req.Phone != nil {
		if !ValidPhoneNumber(*req.Phone) {
			return ErrInvalidPhoneNumber
		}
		if phoneExists, err := u.repository.ExistsCourierByPhone(ctx, *req.Phone); err != nil {
			return err
		} else if phoneExists {
			return ErrPhoneNumberExists
		}
	}

	update := courierUpdateRequestToCourierDB(*req)
	if err := u.repository.UpdateCourier(ctx, &update); err != nil {
		if errors.Is(err, repository.ErrCourierNotFound) {
			return ErrCourierNotFound
		}
		return err
	}
	return nil
}

func ValidPhoneNumber(phone string) bool {
	return regexp.MustCompile(core.PhoneRegex).MatchString(phone)
}
