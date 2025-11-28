package usecase

import (
	"context"
	"courier-service/internal/core"
	"courier-service/internal/model"
	"courier-service/internal/repository"
	"errors"
	"regexp"
)

type CourierUseCase struct {
	repository сourierRepository
}

func NewCourierUseCase(repository сourierRepository) *CourierUseCase {
	return &CourierUseCase{repository: repository}
}

func (u CourierUseCase) GetById(ctx context.Context, id int64) (*model.Courier, error) {
	courierDB, err := u.repository.GetById(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrCourierNotFound) {
			return nil, ErrCourierNotFound
		}
		return nil, err
	}

	courier := courierDBToCourier(*courierDB)
	return &courier, nil
}

func (u CourierUseCase) GetAll(ctx context.Context) ([]model.Courier, error) {
	couriersDB, err := u.repository.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	couriers := make([]model.Courier, len(couriersDB))
	for i, c := range couriersDB {
		couriers[i] = courierDBToCourier(c)
	}

	return couriers, nil
}

func (u CourierUseCase) Create(ctx context.Context, req *model.CourierCreateRequest) (int64, error) {
	if req.Name == "" || req.Phone == "" || req.Status == "" {
		return 0, ErrInvalidCreate
	}

	if !ValidPhoneNumber(req.Phone) {
		return 0, ErrInvalidPhoneNumber
	}

	if phoneExists, err := u.repository.ExistsByPhone(ctx, req.Phone); phoneExists {
		return 0, ErrPhoneNumberExists
	} else if err != nil {
		return 0, err
	}

	courierDB := courierCreateRequestToCourierDB(*req)
	return u.repository.Create(ctx, &courierDB)
}

func (u CourierUseCase) Update(ctx context.Context, req *model.CourierUpdateRequest) error {
	if req.Name == nil && req.Phone == nil && req.Status == nil {
		return ErrInvalidUpdate
	}

	if req.Phone != nil {
		if !ValidPhoneNumber(*req.Phone) {
			return ErrInvalidPhoneNumber
		}
		if phoneExists, err := u.repository.ExistsByPhone(ctx, *req.Phone); err != nil {
			return err
		} else if phoneExists {
			return ErrPhoneNumberExists
		}
	}

	update := courierUpdateRequestToCourierDB(*req)
	if err := u.repository.Update(ctx, &update); err != nil {
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
