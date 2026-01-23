//go:build integration
// +build integration

package delivery_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"courier-service/internal/model"
	integration "courier-service/internal/persistence/database/integration"
	courierstorage "courier-service/internal/repository/courier"
	deliverystorage "courier-service/internal/repository/delivery"
)

type DeliveryTestSuite struct {
	suite.Suite
	ctx          context.Context
	pool         *pgxpool.Pool
	deliveryRepo *deliverystorage.DeliveryRepository
	courierRepo  *courierstorage.CourierRepository
	pgContainer  *postgres.PostgresContainer
	ctrl         *gomock.Controller
	mockLogger   *Mocklogger
}

func (s *DeliveryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	container, connStr, err := integration.TestWithMigrations()
	s.Require().NoError(err)
	s.pgContainer = container

	pool, err := pgxpool.New(s.ctx, connStr)
	s.Require().NoError(err)
	s.pool = pool
	s.deliveryRepo = deliverystorage.NewDeliveryRepository(s.pool)
}

func (s *DeliveryTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockLogger = NewMocklogger(s.ctrl)

	// Allow any logging calls during tests
	s.mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Debugw(gomock.Any(), gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	s.mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

	s.courierRepo = courierstorage.NewCourierRepository(s.pool, s.mockLogger)

	err := integration.TruncateAll(s.ctx, s.pool)
	s.Require().NoError(err)
}

func (s *DeliveryTestSuite) TearDownTest() {
	if s.ctrl != nil {
		s.ctrl.Finish()
	}
}

func TestDeliveryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DeliveryTestSuite))
}

func (s *DeliveryTestSuite) createTestCourier(name, phone string, transportType model.CourierTransportType) int64 {
	courier := model.Courier{
		Name:          name,
		Phone:         phone,
		Status:        model.CourierStatusAvailable,
		TransportType: transportType,
	}
	id, err := s.courierRepo.CreateCourier(context.Background(), courier)
	if err != nil {
		s.T().Fatalf("Failed to create test courier: %v", err)
	}
	return id
}

func (s *DeliveryTestSuite) TestCreate() {
	ctx := context.Background()

	type expectationsFn func(result model.Delivery, err error)

	tests := []struct {
		name         string
		setup        func() model.Delivery
		before       func(d model.Delivery)
		expectations expectationsFn
	}{
		{
			name: "success",
			setup: func() model.Delivery {
				courierID := s.createTestCourier("Test Courier", "+79990000001", "car")

				orderID := uuid.New().String()
				now := time.Now()
				deadline := now.Add(2 * time.Hour)

				return model.Delivery{
					CourierID:  courierID,
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   deadline,
				}
			},
			before: nil,
			expectations: func(result model.Delivery, err error) {
				s.Require().NoError(err)
				s.Require().NotNil(result)
				s.Greater(result.ID, int64(0))
				s.Equal(result.CourierID, result.CourierID)
				s.Equal(result.OrderID, result.OrderID)
			},
		},
		{
			name: "duplicate order id",
			setup: func() model.Delivery {
				// ВАЖНО: другой телефон
				courierID := s.createTestCourier("Test Courier", "+79990000002", model.TransportTypeScooter)

				orderID := uuid.New().String()
				now := time.Now()

				return model.Delivery{
					CourierID:  courierID,
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   now.Add(2 * time.Hour),
				}
			},
			before: func(d model.Delivery) {
				_, err := s.deliveryRepo.CreateDelivery(ctx, d)
				s.Require().NoError(err)
			},
			expectations: func(result model.Delivery, err error) {
				s.ErrorIs(err, deliverystorage.ErrOrderIDExists)
				s.Equal(model.Delivery{}, result)
			},
		},
		{
			name: "foreign key violation",
			setup: func() model.Delivery {
				orderID := uuid.New().String()
				now := time.Now()

				return model.Delivery{
					CourierID:  999999, // несуществующий курьер
					OrderID:    orderID,
					AssignedAt: now,
					Deadline:   now.Add(2 * time.Hour),
				}
			},
			before: nil,
			expectations: func(result model.Delivery, err error) {
				s.Error(err)
				s.Equal(model.Delivery{}, result)
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

func (s *DeliveryTestSuite) TestCouriersDelivery() {
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

				deliveryDB := model.Delivery{
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

				s.ErrorIs(err, deliverystorage.ErrOrderIDNotFound)
				s.Equal(model.Delivery{}, result)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.test()
		})
	}
}

func (s *DeliveryTestSuite) TestDelete() {
	ctx := context.Background()

	tests := []struct {
		name  string
		setup func() string
		check func(orderID string)
	}{
		{
			name: "success",
			setup: func() string {
				courierID := s.createTestCourier("Test Courier", "+79991234567", "car")

				orderID := uuid.New().String()
				now := time.Now()

				deliveryDB := model.Delivery{
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
				s.ErrorIs(err, deliverystorage.ErrOrderIDNotFound)
			},
		},
		{
			name: "not found",
			setup: func() string {
				return uuid.New().String()
			},
			check: func(orderID string) {
				err := s.deliveryRepo.DeleteDelivery(ctx, orderID)
				s.ErrorIs(err, deliverystorage.ErrOrderIDNotFound)
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

func (s *DeliveryTestSuite) TestMultipleDeliveries() {
	courierID1 := s.createTestCourier("Courier 1", "+79991234567", model.TransportTypeCar)
	courierID2 := s.createTestCourier("Courier 2", "+79991234568", model.TransportTypeScooter)

	orderID1 := uuid.New().String()
	orderID2 := uuid.New().String()
	now := time.Now()

	delivery1 := model.Delivery{
		CourierID:  courierID1,
		OrderID:    orderID1,
		AssignedAt: now,
		Deadline:   now.Add(2 * time.Hour),
	}

	delivery2 := model.Delivery{
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
