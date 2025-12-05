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


func (s *DeliveryRepositoryTestSuite) TestCreate() {
	ctx := context.Background()

	type expectationsFn func(result *model.Delivery, err error)

	tests := []struct {
		name         string
		setup        func() *model.DeliveryDB
		before       func(d *model.DeliveryDB)
		expectations expectationsFn
	}{
		{
			name: "success",
			setup: func() *model.DeliveryDB {
				courierID := s.createTestCourier("Test Courier", "+79990000001", "car")

				orderID := uuid.New().String()
				now := time.Now()
				deadline := now.Add(2 * time.Hour)

				return &model.DeliveryDB{
					CourierID:  courierID,
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   deadline,
				}
			},
			before: nil,
			expectations: func(result *model.Delivery, err error) {
				s.Require().NoError(err)
				s.Require().NotNil(result)
				s.Greater(result.ID, int64(0))
				s.Equal(result.CourierID, result.CourierID)
				s.Equal(result.OrderID, result.OrderID)
			},
		},
		{
			name: "duplicate order id",
			setup: func() *model.DeliveryDB {
				// ВАЖНО: другой телефон
				courierID := s.createTestCourier("Test Courier", "+79990000002", "car")

				orderID := uuid.New().String()
				now := time.Now()

				return &model.DeliveryDB{
					CourierID:  courierID,
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   now.Add(2 * time.Hour),
				}
			},
			before: func(d *model.DeliveryDB) {
				_, err := s.deliveryRepo.CreateDelivery(ctx, d)
				s.Require().NoError(err)
			},
			expectations: func(result *model.Delivery, err error) {
				s.ErrorIs(err, ErrOrderIDExists)
				s.Nil(result)
			},
		},
		{
			name: "foreign key violation",
			setup: func() *model.DeliveryDB {
				orderID := uuid.New().String()
				now := time.Now()

				return &model.DeliveryDB{
					CourierID:  999999, // несуществующий курьер
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   now.Add(2 * time.Hour),
				}
			},
			before: nil,
			expectations: func(result *model.Delivery, err error) {
				s.Error(err)
				s.Nil(result)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			d := tt.setup()

			if tt.before != nil {
				tt.before(d)
			}

			result, err := s.deliveryRepo.CreateDelivery(ctx, d)

			tt.expectations(result, err)
		})
	}
}

func (s *DeliveryRepositoryTestSuite) TestCouriersDelivery() {
	tests := []struct {
		name string
		test func()
	}{
		{
			name: "success",
			test: func() {
				ctx := context.Background()

				courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

				orderID := uuid.New().String()
				now := time.Now()

				deliveryDB := &model.DeliveryDB{
					CourierID:  courierID,
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   now.Add(2 * time.Hour),
				}

				_, err := s.deliveryRepo.CreateDelivery(ctx, deliveryDB)
				s.Require().NoError(err)

				result, err := s.deliveryRepo.CouriersDelivery(ctx, orderID)

				s.Require().NoError(err)
				s.Require().NotNil(result)
				s.Equal(courierID, result.CourierID)
				s.Equal(orderID, result.OrderID)
			},
		},
		{
			name: "not found",
			test: func() {
				ctx := context.Background()

				nonExistentOrderID := uuid.New().String()

				result, err := s.deliveryRepo.CouriersDelivery(ctx, nonExistentOrderID)

				s.ErrorIs(err, ErrOrderIDNotFound)
				s.Nil(result)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.test()
		})
	}
}

func (s *DeliveryRepositoryTestSuite) TestDelete() {
	ctx := context.Background()

	tests := []struct {
		name   string
		setup  func() string
		check  func(orderID string)
	}{
		{
			name: "success",
			setup: func() string {
				courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

				orderID := uuid.New().String()
				now := time.Now()

				deliveryDB := &model.DeliveryDB{
					CourierID:  courierID,
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   now.Add(2 * time.Hour),
				}

				_, err := s.deliveryRepo.CreateDelivery(ctx, deliveryDB)
				s.Require().NoError(err)

				return orderID
			},
			check: func(orderID string) {
				err := s.deliveryRepo.DeleteDelivery(ctx, orderID)
				s.Require().NoError(err)

				_, err = s.deliveryRepo.CouriersDelivery(ctx, orderID)
				s.ErrorIs(err, ErrOrderIDNotFound)
			},
		},
		{
			name: "not found",
			setup: func() string {
				return uuid.New().String()
			},
			check: func(orderID string) {
				err := s.deliveryRepo.DeleteDelivery(ctx, orderID)
				s.ErrorIs(err, ErrOrderIDNotFound)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			orderID := tt.setup()
			tt.check(orderID)
		})
	}
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
