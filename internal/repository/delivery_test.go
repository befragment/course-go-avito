package repository

import (
	"context"
	"courier-service/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type DeliveryRepositoryTestSuite struct {
	RepositoryTestSuite
}

func TestDeliveryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryRepositoryTestSuite))
}

func (s *DeliveryRepositoryTestSuite) TestCreate_Success() {
	courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

	orderID := uuid.New().String()
	now := time.Now()
	deadline := now.Add(2 * time.Hour)

	deliveryDB := &model.DeliveryDB{
		CourierID:  courierID,
		OrderID:    orderID,
		AssignedAt: now,
		Deadline:   deadline,
	}

	result, err := s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Greater(result.ID, int64(0))
	s.Equal(courierID, result.CourierID)
	s.Equal(orderID, result.OrderID)
	s.WithinDuration(deadline, result.Deadline, time.Second)
}

func (s *DeliveryRepositoryTestSuite) TestCreate_DuplicateOrderID() {
	courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

	orderID := uuid.New().String()
	now := time.Now()

	deliveryDB := &model.DeliveryDB{
		CourierID:  courierID,
		OrderID:    orderID,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	_, err := s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)
	s.Require().NoError(err)

	_, err = s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)

	s.ErrorIs(err, ErrOrderIDExists)
}

func (s *DeliveryRepositoryTestSuite) TestCreate_ForeignKeyViolation() {
	orderID := uuid.New().String()
	now := time.Now()

	deliveryDB := &model.DeliveryDB{
		CourierID:  999999, // Non-existent courier
		OrderID:    orderID,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	_, err := s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)

	s.Error(err)
}

func (s *DeliveryRepositoryTestSuite) TestCouriersDelivery_Success() {
	courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

	orderID := uuid.New().String()
	now := time.Now()

	deliveryDB := &model.DeliveryDB{
		CourierID:  courierID,
		OrderID:    orderID,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	_, err := s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)
	s.Require().NoError(err)

	result, err := s.deliveryRepo.CouriersDelivery(context.Background(), orderID)

	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Equal(courierID, result.CourierID)
	s.Equal(orderID, result.OrderID)
}

func (s *DeliveryRepositoryTestSuite) TestCouriersDelivery_NotFound() {
	nonExistentOrderID := uuid.New().String()

	result, err := s.deliveryRepo.CouriersDelivery(context.Background(), nonExistentOrderID)

	s.ErrorIs(err, ErrOrderIDNotFound)
	s.Nil(result)
}

func (s *DeliveryRepositoryTestSuite) TestDelete_Success() {
	courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

	orderID := uuid.New().String()
	now := time.Now()

	deliveryDB := &model.DeliveryDB{
		CourierID:  courierID,
		OrderID:    orderID,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	_, err := s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)
	s.Require().NoError(err)

	err = s.deliveryRepo.DeleteDelivery(context.Background(), orderID)

	s.Require().NoError(err)

	_, err = s.deliveryRepo.CouriersDelivery(context.Background(), orderID)
	s.ErrorIs(err, ErrOrderIDNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestDelete_NotFound() {
	nonExistentOrderID := uuid.New().String()

	err := s.deliveryRepo.DeleteDelivery(context.Background(), nonExistentOrderID)

	s.ErrorIs(err, ErrOrderIDNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestDelete_MultipleTimes() {
	courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

	orderID := uuid.New().String()
	now := time.Now()

	deliveryDB := &model.DeliveryDB{
		CourierID:  courierID,
		OrderID:    orderID,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	_, err := s.deliveryRepo.CreateDelivery(context.Background(), deliveryDB)
	s.Require().NoError(err)

	err = s.deliveryRepo.DeleteDelivery(context.Background(), orderID)
	s.Require().NoError(err)

	err = s.deliveryRepo.DeleteDelivery(context.Background(), orderID)
	s.ErrorIs(err, ErrOrderIDNotFound)
}

func (s *DeliveryRepositoryTestSuite) TestMultipleDeliveries() {
	courierID1 := s.createTestCourier("Courier 1", "+79991234567", "car")
	courierID2 := s.createTestCourier("Courier 2", "+79991234568", "bike")

	orderID1 := uuid.New().String()
	orderID2 := uuid.New().String()
	now := time.Now()

	delivery1 := &model.DeliveryDB{
		CourierID:  courierID1,
		OrderID:    orderID1,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	delivery2 := &model.DeliveryDB{
		CourierID:  courierID2,
		OrderID:    orderID2,
		AssignedAt: now,
		Deadline:   now.Add(3 * time.Hour),
	}

	_, err := s.deliveryRepo.CreateDelivery(context.Background(), delivery1)
	s.Require().NoError(err)

	_, err = s.deliveryRepo.CreateDelivery(context.Background(), delivery2)
	s.Require().NoError(err)

	result1, err := s.deliveryRepo.CouriersDelivery(context.Background(), orderID1)
	s.Require().NoError(err)
	s.Equal(courierID1, result1.CourierID)

	result2, err := s.deliveryRepo.CouriersDelivery(context.Background(), orderID2)
	s.Require().NoError(err)
	s.Equal(courierID2, result2.CourierID)
}
