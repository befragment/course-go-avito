package repository

import (
	"context"
	"courier-service/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CourierRepositoryTestSuite struct {
	RepositoryTestSuite
}

func TestCourierRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CourierRepositoryTestSuite))
}

func (s *CourierRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()
	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	id, err := s.courierRepo.CreateCourier(ctx, courier)

	s.Require().NoError(err)
	s.Greater(id, int64(0))
}

func (s *CourierRepositoryTestSuite) TestCreate_DuplicatePhone() {
	ctx := context.Background()
	phone := "+79991234567"

	courier1 := &model.CourierDB{
		Name:          "John Doe",
		Phone:         phone,
		Status:        "available",
		TransportType: "car",
	}
	_, err := s.courierRepo.CreateCourier(ctx, courier1)
	s.Require().NoError(err)

	courier2 := &model.CourierDB{
		Name:          "Jane Doe",
		Phone:         phone,
		Status:        "available",
		TransportType: "bike",
	}
	_, err = s.courierRepo.CreateCourier(ctx, courier2)

	s.Require().Error(err)
	s.ErrorIs(err, ErrPhoneNumberExists)
}

func (s *CourierRepositoryTestSuite) TestGetById_Success() {
	ctx := context.Background()
	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	result, err := s.courierRepo.GetCourierById(ctx, id)

	s.Require().NoError(err)
	s.Equal(id, result.ID)
	s.Equal("John Doe", result.Name)
	s.Equal("+79991234567", result.Phone)
	s.Equal("available", result.Status)
	s.Equal("car", result.TransportType)
}

func (s *CourierRepositoryTestSuite) TestGetById_NotFound() {
	ctx := context.Background()

	result, err := s.courierRepo.GetCourierById(ctx, 99999)

	s.Require().Error(err)
	s.Nil(result)
	s.ErrorIs(err, ErrCourierNotFound)
}

func (s *CourierRepositoryTestSuite) TestGetAll_Success() {
	ctx := context.Background()

	couriers := []*model.CourierDB{
		{Name: "John", Phone: "+79991234567", Status: "available", TransportType: "car"},
		{Name: "Jane", Phone: "+79991234568", Status: "available", TransportType: "scooter"},
		{Name: "Bob", Phone: "+79991234569", Status: "available", TransportType: "on_foot"},
	}

	for _, c := range couriers {
		_, err := s.courierRepo.CreateCourier(ctx, c)
		s.Require().NoError(err)
	}

	result, err := s.courierRepo.GetAllCouriers(ctx)

	s.Require().NoError(err)
	s.Len(result, 3)
	s.Equal("John", result[0].Name)
	s.Equal("Jane", result[1].Name)
	s.Equal("Bob", result[2].Name)
}

func (s *CourierRepositoryTestSuite) TestGetAll_Empty() {
	ctx := context.Background()

	result, err := s.courierRepo.GetAllCouriers(ctx)

	s.Require().NoError(err)
	s.Empty(result)
}

func (s *CourierRepositoryTestSuite) TestUpdate_Success() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	updated := &model.CourierDB{
		ID:            id,
		Name:          "Jane Doe",
		Status:        "busy",
		TransportType: "bike",
	}

	err = s.courierRepo.UpdateCourier(ctx, updated)
	s.Require().NoError(err)

	result, err := s.courierRepo.GetCourierById(ctx, id)
	s.Require().NoError(err)
	s.Equal("Jane Doe", result.Name)
	s.Equal("busy", result.Status)
	s.Equal("bike", result.TransportType)
}

func (s *CourierRepositoryTestSuite) TestUpdate_NotFound() {
	ctx := context.Background()

	courier := &model.CourierDB{
		ID:            99999,
		Name:          "John Doe",
		Status:        "available",
		TransportType: "car",
	}

	err := s.courierRepo.UpdateCourier(ctx, courier)

	s.Require().Error(err)
	s.ErrorIs(err, ErrCourierNotFound)
}

func (s *CourierRepositoryTestSuite) TestUpdate_NothingToUpdate() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	emptyUpdate := &model.CourierDB{
		ID: id,
	}

	err = s.courierRepo.UpdateCourier(ctx, emptyUpdate)

	s.Require().Error(err)
	s.ErrorIs(err, ErrNothingToUpdate)
}

func (s *CourierRepositoryTestSuite) TestExistsByPhone_True() {
	ctx := context.Background()
	phone := "+79991234567"

	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         phone,
		Status:        "available",
		TransportType: "car",
	}
	_, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	exists, err := s.courierRepo.ExistsCourierByPhone(ctx, phone)

	s.Require().NoError(err)
	s.True(exists)
}

func (s *CourierRepositoryTestSuite) TestExistsByPhone_False() {
	ctx := context.Background()

	exists, err := s.courierRepo.ExistsCourierByPhone(ctx, "+79999999999")

	s.Require().NoError(err)
	s.False(exists)
}

func (s *CourierRepositoryTestSuite) TestFindAvailable_Success() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	_, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	result, err := s.courierRepo.FindAvailableCourier(ctx)

	s.Require().NoError(err)
	s.NotNil(result)
	s.Equal("available", result.Status)
}

func (s *CourierRepositoryTestSuite) TestFindAvailable_AllBusy() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	busyUpdate := &model.CourierDB{
		ID:     id,
		Status: "busy",
	}
	err = s.courierRepo.UpdateCourier(ctx, busyUpdate)
	s.Require().NoError(err)

	result, err := s.courierRepo.FindAvailableCourier(ctx)

	s.Require().Error(err)
	s.Nil(result)
	s.ErrorIs(err, ErrCouriersBusy)
}

func (s *CourierRepositoryTestSuite) TestFindAvailable_SelectsCourierWithFewestDeliveries() {
	ctx := context.Background()

	courier1 := &model.CourierDB{
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id1, err := s.courierRepo.CreateCourier(ctx, courier1)
	s.Require().NoError(err)

	courier2 := &model.CourierDB{
		Name:          "Jane",
		Phone:         "+79991234568",
		Status:        "available",
		TransportType: "car",
	}
	id2, err := s.courierRepo.CreateCourier(ctx, courier2)
	s.Require().NoError(err)

	orderID := uuid.New().String()
	_, err = s.pool.Exec(ctx,
		"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
		id1, orderID, time.Now(), time.Now().Add(time.Hour))
	s.Require().NoError(err)

	result, err := s.courierRepo.FindAvailableCourier(ctx)

	s.Require().NoError(err)
	s.NotNil(result)
	s.Equal(id2, result.ID)
}

func (s *CourierRepositoryTestSuite) TestFreeCouriers_Success() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	busyUpdate := &model.CourierDB{
		ID:     id,
		Status: "busy",
	}
	err = s.courierRepo.UpdateCourier(ctx, busyUpdate)
	s.Require().NoError(err)

	pastTime := time.Now().Add(-1 * time.Hour)
	orderID := uuid.New().String()
	_, err = s.pool.Exec(ctx,
		"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
		id, orderID, pastTime, pastTime)
	s.Require().NoError(err)

	err = s.courierRepo.FreeCouriersWithInterval(ctx)
	s.Require().NoError(err)

	updated, err := s.courierRepo.GetCourierById(ctx, id)
	s.Require().NoError(err)
	s.Equal("available", updated.Status)
}

func (s *CourierRepositoryTestSuite) TestFreeCouriers_NoExpiredDeliveries() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	busyUpdate := &model.CourierDB{
		ID:     id,
		Status: "busy",
	}
	err = s.courierRepo.UpdateCourier(ctx, busyUpdate)
	s.Require().NoError(err)

	futureTime := time.Now().Add(1 * time.Hour)
	orderID := uuid.New().String()
	_, err = s.pool.Exec(ctx,
		"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
		id, orderID, time.Now(), futureTime)
	s.Require().NoError(err)

	err = s.courierRepo.FreeCouriersWithInterval(ctx)
	s.Require().NoError(err)

	updated, err := s.courierRepo.GetCourierById(ctx, id)
	s.Require().NoError(err)
	s.Equal("busy", updated.Status)
}

func (s *CourierRepositoryTestSuite) TestFreeCouriers_OnlyFreesLastExpiredDelivery() {
	ctx := context.Background()

	courier := &model.CourierDB{
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.courierRepo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	busyUpdate := &model.CourierDB{
		ID:     id,
		Status: "busy",
	}
	err = s.courierRepo.UpdateCourier(ctx, busyUpdate)
	s.Require().NoError(err)

	pastTime := time.Now().Add(-2 * time.Hour)
	orderID1 := uuid.New().String()
	_, err = s.pool.Exec(ctx,
		"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
		id, orderID1, pastTime, pastTime)
	s.Require().NoError(err)

	futureTime := time.Now().Add(1 * time.Hour)
	orderID2 := uuid.New().String()
	_, err = s.pool.Exec(ctx,
		"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
		id, orderID2, time.Now(), futureTime)
	s.Require().NoError(err)

	err = s.courierRepo.FreeCouriersWithInterval(ctx)
	s.Require().NoError(err)

	updated, err := s.courierRepo.GetCourierById(ctx, id)
	s.Require().NoError(err)
	s.Equal("busy", updated.Status)
}
