//go:build integration
// +build integration

package courier_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/golang/mock/gomock"

	"courier-service/internal/model"
	"courier-service/internal/persistence/database/integration"
	courierstorage "courier-service/internal/repository/courier"
)

type CourierTestSuite struct {
	suite.Suite
	ctx         context.Context
	pool        *pgxpool.Pool
	repo        *courierstorage.CourierRepository
	pgContainer *postgres.PostgresContainer
	ctrl        *gomock.Controller
	mockLogger  *Mocklogger
}

func (s *CourierTestSuite) SetupSuite() {
	s.ctx = context.Background()

	container, connStr, err := integration.TestWithMigrations()
	s.Require().NoError(err)
	s.pgContainer = container

	pool, err := pgxpool.New(s.ctx, connStr)
	s.Require().NoError(err)
	s.pool = pool
}

func (s *CourierTestSuite) SetupTest() {
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

	s.repo = courierstorage.NewCourierRepository(s.pool, s.mockLogger)

	err := integration.TruncateAll(s.ctx, s.pool)
	s.Require().NoError(err)
}

func (s *CourierTestSuite) TearDownTest() {
	if s.ctrl != nil {
		s.ctrl.Finish()
	}
}

func TestCourierRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CourierTestSuite))
}

func (s *CourierTestSuite) TestCreate() {
	ctx := context.Background()

	tests := []struct {
		name string
		test func()
	}{
		{
			name: "success",
			test: func() {
				courier := model.Courier{
					Name:          "John Doe",
					Phone:         "+79990000001",
					Status:        "available",
					TransportType: "car",
				}

				id, err := s.repo.CreateCourier(ctx, courier)

				s.Require().NoError(err)
				s.Greater(id, int64(0))
			},
		},
		{
			name: "duplicate phone",
			test: func() {
				phone := "+79990000002"

				courier1 := model.Courier{
					Name:          "John Doe",
					Phone:         phone,
					Status:        "available",
					TransportType: "car",
				}
				_, err := s.repo.CreateCourier(ctx, courier1)
				s.Require().NoError(err)

				courier2 := model.Courier{
					Name:          "Jane Doe",
					Phone:         phone,
					Status:        "available",
					TransportType: "bike",
				}
				_, err = s.repo.CreateCourier(ctx, courier2)

				s.Require().Error(err)
				s.ErrorIs(err, courierstorage.ErrPhoneNumberExists)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.test()
		})
	}
}

func (s *CourierTestSuite) TestGetById_Success() {
	ctx := context.Background()
	courier := model.Courier{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        model.CourierStatusAvailable,
		TransportType: "car",
	}

	id, err := s.repo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	result, err := s.repo.GetCourierById(ctx, id)

	s.Require().NoError(err)
	s.Equal(id, result.ID)
	s.Equal("John Doe", result.Name)
	s.Equal("+79991234567", result.Phone)
	s.Equal(model.CourierStatusAvailable, result.Status)
	s.Equal(model.TransportTypeCar, result.TransportType)
}

func (s *CourierTestSuite) TestGetById_NotFound() {
	ctx := context.Background()

	result, err := s.repo.GetCourierById(ctx, 99999)

	s.Require().Error(err)
	s.Equal(model.Courier{}, result)
	s.ErrorIs(err, courierstorage.ErrCourierNotFound)
}

func (s *CourierTestSuite) TestGetAll_Success() {
	ctx := context.Background()

	couriers := []model.Courier{
		{Name: "John", Phone: "+79991234567", Status: model.CourierStatusAvailable, TransportType: "car"},
		{Name: "Jane", Phone: "+79991234568", Status: model.CourierStatusAvailable, TransportType: "scooter"},
		{Name: "Bob", Phone: "+79991234569", Status: model.CourierStatusAvailable, TransportType: "on_foot"},
	}

	for _, c := range couriers {
		_, err := s.repo.CreateCourier(ctx, c)
		s.Require().NoError(err)
	}

	result, err := s.repo.GetAllCouriers(ctx)

	s.Require().NoError(err)
	s.Len(result, 3)
	s.Equal("John", result[0].Name)
	s.Equal("Jane", result[1].Name)
	s.Equal("Bob", result[2].Name)
}

func (s *CourierTestSuite) TestGetAll_Empty() {
	ctx := context.Background()

	result, err := s.repo.GetAllCouriers(ctx)

	s.Require().NoError(err)
	s.Empty(result)
}

func (s *CourierTestSuite) TestUpdate_Success() {
	ctx := context.Background()

	courier := model.Courier{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        model.CourierStatusAvailable,
		TransportType: "car",
	}
	id, err := s.repo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	updated := model.Courier{
		ID:            id,
		Name:          "Jane Doe",
		Status:        model.CourierStatusBusy,
		TransportType: "bike",
	}

	err = s.repo.UpdateCourier(ctx, updated)
	s.Require().NoError(err)

	result, err := s.repo.GetCourierById(ctx, id)
	s.Require().NoError(err)
	s.Equal("Jane Doe", result.Name)
	s.Equal(model.CourierStatusBusy, result.Status)
	s.Equal(model.CourierTransportType("bike"), result.TransportType)
}

func (s *CourierTestSuite) TestUpdate_NotFound() {
	ctx := context.Background()

	courier := model.Courier{
		ID:            99999,
		Name:          "John Doe",
		Status:        "available",
		TransportType: "car",
	}

	err := s.repo.UpdateCourier(ctx, courier)

	s.Require().Error(err)
	s.ErrorIs(err, courierstorage.ErrCourierNotFound)
}

func (s *CourierTestSuite) TestUpdate_NothingToUpdate() {
	ctx := context.Background()

	courier := model.Courier{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	id, err := s.repo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	emptyUpdate := model.Courier{
		ID: id,
	}

	err = s.repo.UpdateCourier(ctx, emptyUpdate)

	s.Require().Error(err)
	s.ErrorIs(err, courierstorage.ErrNothingToUpdate)
}

func (s *CourierTestSuite) TestExistsByPhone_True() {
	ctx := context.Background()
	phone := "+79991234567"

	courier := model.Courier{
		Name:          "John Doe",
		Phone:         phone,
		Status:        model.CourierStatusAvailable,
		TransportType: "car",
	}
	_, err := s.repo.CreateCourier(ctx, courier)
	s.Require().NoError(err)

	exists, err := s.repo.ExistsCourierByPhone(ctx, phone)

	s.Require().NoError(err)
	s.True(exists)
}

func (s *CourierTestSuite) TestExistsByPhone_False() {
	ctx := context.Background()

	exists, err := s.repo.ExistsCourierByPhone(ctx, "+79999999999")

	s.Require().NoError(err)
	s.False(exists)
}

func (s *CourierTestSuite) TestFindAvailable() {
	ctx := context.Background()

	tests := []struct {
		name string
		test func()
	}{
		{
			name: "success",
			test: func() {
				courier := model.Courier{
					Name:          "John Doe",
					Phone:         "+79990000001",
					Status:        model.CourierStatusAvailable,
					TransportType: "car",
				}
				_, err := s.repo.CreateCourier(ctx, courier)
				s.Require().NoError(err)

				result, err := s.repo.FindAvailableCourier(ctx)

				s.Require().NoError(err)
				s.NotNil(result)
				s.Equal(model.CourierStatusAvailable, result.Status)
			},
		},
		{
			name: "all_busy",
			test: func() {
				courier := model.Courier{
					Name:          "John Doe",
					Phone:         "+79990000002",
					Status:        model.CourierStatusAvailable,
					TransportType: "car",
				}
				id, err := s.repo.CreateCourier(ctx, courier)
				s.Require().NoError(err)

				busyUpdate := model.Courier{
					ID:     id,
					Status: model.CourierStatusBusy,
				}
				err = s.repo.UpdateCourier(ctx, busyUpdate)
				s.Require().NoError(err)

				result, err := s.repo.FindAvailableCourier(ctx)

				s.Require().Error(err)
				s.Equal(model.Courier{}, result)
				s.ErrorIs(err, courierstorage.ErrCouriersBusy)
			},
		},
		{
			name: "selects_courier_with_fewest_deliveries",
			test: func() {
				// courier1 — с одной доставкой
				courier1 := model.Courier{
					Name:          "John",
					Phone:         "+79990000003",
					Status:        model.CourierStatusAvailable,
					TransportType: "car",
				}
				id1, err := s.repo.CreateCourier(ctx, courier1)
				s.Require().NoError(err)

				// courier2 — без доставок
				courier2 := model.Courier{
					Name:          "Jane",
					Phone:         "+79990000004",
					Status:        model.CourierStatusAvailable,
					TransportType: "car",
				}
				id2, err := s.repo.CreateCourier(ctx, courier2)
				s.Require().NoError(err)

				orderID := uuid.New().String()
				_, err = s.pool.Exec(ctx,
					"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
					id1, orderID, time.Now(), time.Now().Add(time.Hour))
				s.Require().NoError(err)

				result, err := s.repo.FindAvailableCourier(ctx)

				s.Require().NoError(err)
				s.NotNil(result)
				s.Equal(id2, result.ID)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Чистим состояние БД перед каждым сабтестом,
			// чтобы кейсы не влияли друг на друга.
			_, err := s.pool.Exec(ctx, `
				TRUNCATE TABLE delivery, couriers RESTART IDENTITY CASCADE;
			`)
			s.Require().NoError(err)

			tt.test()
		})
	}
}

func (s *CourierTestSuite) TestFreeCouriers() {
	ctx := context.Background()

	tests := []struct {
		name            string
		phone           string
		setupDeliveries func(id int64)
		expectedStatus  model.CourierStatus
	}{
		{
			name:  "success_expired_delivery_frees_courier",
			phone: "+79990000001",
			setupDeliveries: func(id int64) {
				pastTime := time.Now().Add(-1 * time.Hour)
				orderID := uuid.New().String()

				_, err := s.pool.Exec(ctx,
					"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
					id, orderID, pastTime, pastTime)
				s.Require().NoError(err)
			},
			expectedStatus: model.CourierStatusAvailable,
		},
		{
			name:  "no_expired_deliveries_courier_stays_busy",
			phone: "+79990000002",
			setupDeliveries: func(id int64) {
				futureTime := time.Now().Add(1 * time.Hour)
				orderID := uuid.New().String()

				_, err := s.pool.Exec(ctx,
					"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
					id, orderID, time.Now(), futureTime)
				s.Require().NoError(err)
			},
			expectedStatus: model.CourierStatusBusy,
		},
		{
			name:  "only_frees_when_all_deliveries_expired",
			phone: "+79990000003",
			setupDeliveries: func(id int64) {
				pastTime := time.Now().Add(-2 * time.Hour)
				orderID1 := uuid.New().String()

				_, err := s.pool.Exec(ctx,
					"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
					id, orderID1, pastTime, pastTime)
				s.Require().NoError(err)

				futureTime := time.Now().Add(1 * time.Hour)
				orderID2 := uuid.New().String()

				_, err = s.pool.Exec(ctx,
					"INSERT INTO delivery (courier_id, order_id, assigned_at, deadline) VALUES ($1, $2, $3, $4)",
					id, orderID2, time.Now(), futureTime)
				s.Require().NoError(err)
			},
			expectedStatus: model.CourierStatusBusy,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// создаём курьера
			courier := model.Courier{
				Name:          "John",
				Phone:         tt.phone,
				Status:        model.CourierStatus(model.CourierStatusAvailable),
				TransportType: "car",
			}
			id, err := s.repo.CreateCourier(ctx, courier)
			s.Require().NoError(err)

			// помечаем как busy
			busyUpdate := model.Courier{
				ID:     id,
				Status: model.CourierStatusBusy,
			}
			err = s.repo.UpdateCourier(ctx, busyUpdate)
			s.Require().NoError(err)

			// вставляем доставки по сценарию
			tt.setupDeliveries(id)

			// вызываем освобождение
			err = s.repo.FreeCouriersWithInterval(ctx)
			s.Require().NoError(err)

			// проверяем статус
			updated, err := s.repo.GetCourierById(ctx, id)
			s.Require().NoError(err)
			s.Equal(tt.expectedStatus, updated.Status)
		})
	}
}
