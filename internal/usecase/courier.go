package usecase

import (
	"context"
	"errors"
	"regexp"

	"courier-service/internal/core"
	"courier-service/internal/model"
	"courier-service/internal/repository"
)

type CourierUseCase struct {
	repository CourierRepository
}

func NewCourierUseCase(repository CourierRepository) *CourierUseCase {
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

	courier := model.Courier{
		ID:        courierDB.ID,
		Name:      courierDB.Name,
		Phone:     courierDB.Phone,
		Status:    courierDB.Status,
		CreatedAt: courierDB.CreatedAt,
		UpdatedAt: courierDB.UpdatedAt,
	}

	return &courier, nil
}

func (u CourierUseCase) GetAll(ctx context.Context) ([]model.Courier, error) {
	couriersDB, err := u.repository.GetAll(ctx)

	if err != nil {
		return []model.Courier{}, err
	}

	var couriers []model.Courier
	for _, c := range couriersDB {
		couriers = append(couriers, model.Courier{
			ID: 		c.ID,
			Name: 		c.Name,
			Phone: 		c.Phone,
			Status: 	c.Status,
			CreatedAt: 	c.CreatedAt,
			UpdatedAt: 	c.UpdatedAt,
		})
	}

	return couriers, nil
}

func (u CourierUseCase) Create(ctx context.Context, req *model.CourierCreateRequest) (int64, error) {
	if req.Name == "" || req.Phone == "" || req.Status == "" {
		return 0, ErrInvalidCreate
	}

	if !ValidPhoneNumber(&req.Phone) {
		return 0, ErrInvalidPhoneNumber
	}
	
	if phoneExists, err := u.repository.ExistsByPhone(ctx, req.Phone); phoneExists {
		return 0, ErrPhoneNumberExists
	} else if err != nil {
		return 0, err
	}

	courierDB := &model.CourierDB{
		Name: req.Name,
		Phone: req.Phone,
		Status: req.Status,
	}

	id, err := u.repository.Create(ctx, courierDB)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u CourierUseCase) Update(ctx context.Context, req *model.CourierUpdateRequest) error {
	if req.ID == 0 {
		return ErrIdRequired
	}

	if req.Name == nil && req.Phone == nil && req.Status == nil {
		return ErrInvalidUpdate
	}

	if req.Phone != nil && !ValidPhoneNumber(req.Phone) {
		return ErrInvalidPhoneNumber
	}

	if phoneExists, err := u.repository.ExistsByPhone(ctx, *req.Phone); phoneExists {
		return ErrPhoneNumberExists
	} else if err != nil {
		return err
	}

	err := u.repository.Update(ctx, &model.CourierDB{
		ID: req.ID,
		Name: *req.Name,
		Phone: *req.Phone,
		Status: *req.Status,
	})
	
	if err != nil {
		return err
	}

	return nil
}

func ValidPhoneNumber(phone *string) bool {
	return regexp.MustCompile(core.PhoneRegex).MatchString(*phone)
}
