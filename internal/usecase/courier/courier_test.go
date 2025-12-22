package courier_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"courier-service/internal/model"
	courierRepo "courier-service/internal/repository/courier"
	"courier-service/internal/usecase/courier"
	logger "courier-service/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestCourierUseCase_GetById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		courierID    int64
		prepare      func(repo *MockcourierRepository)
		expectations func(t *testing.T, result model.Courier, err error)
	}{
		{
			name:      "success: courier found",
			courierID: 1,
			prepare: func(repo *MockcourierRepository) {
				repo.EXPECT().
					GetCourierById(gomock.Any(), int64(1)).
					Return(model.Courier{
						ID:            1,
						Name:          "John Doe",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)
			},
			expectations: func(t *testing.T, result model.Courier, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int64(1), result.ID)
				assert.Equal(t, "John Doe", result.Name)
				assert.Equal(t, "+79991234567", result.Phone)
			},
		},
		{
			name:      "error: courier not found",
			courierID: 999,
			prepare: func(repo *MockcourierRepository) {
				repo.EXPECT().
					GetCourierById(gomock.Any(), int64(999)).
					Return(model.Courier{}, courierRepo.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, result model.Courier, err error) {
				assert.Error(t, err)
				assert.Equal(t, model.Courier{}, result)
				assert.Equal(t, courier.ErrCourierNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockcourierRepository(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
			logger, err := logger.New(logger.LogLevelInfo)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}
			uc := courier.NewCourierUseCase(mockRepo, mockFactory, logger)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockRepo)
			}

			result, err := uc.GetCourierById(ctx, tc.courierID)

			if tc.expectations != nil {
				tc.expectations(t, result, err)
			}
		})
	}
}

func TestCourierUseCase_GetAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		prepare      func(repo *MockcourierRepository)
		expectations func(t *testing.T, result []model.Courier, err error)
	}{
		{
			name: "success: returns multiple couriers",
			prepare: func(repo *MockcourierRepository) {
				repo.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return([]model.Courier{
						{ID: 1, Name: "John", Phone: "+79991234567", Status: "available", TransportType: "car"},
						{ID: 2, Name: "Jane", Phone: "+79991234568", Status: "busy", TransportType: "scooter"},
					}, nil)
			},
			expectations: func(t *testing.T, result []model.Courier, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 2)
				assert.Equal(t, int64(1), result[0].ID)
				assert.Equal(t, "John", result[0].Name)
				assert.Equal(t, int64(2), result[1].ID)
				assert.Equal(t, "Jane", result[1].Name)
			},
		},
		{
			name: "success: returns empty list",
			prepare: func(repo *MockcourierRepository) {
				repo.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return([]model.Courier{}, nil)
			},
			expectations: func(t *testing.T, result []model.Courier, err error) {
				assert.NoError(t, err)
				assert.Empty(t, result)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockcourierRepository(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
			logger, err := logger.New(logger.LogLevelInfo)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}
			uc := courier.NewCourierUseCase(mockRepo, mockFactory, logger)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockRepo)
			}

			result, err := uc.GetAllCouriers(ctx)

			if tc.expectations != nil {
				tc.expectations(t, result, err)
			}
		})
	}
}

func TestCourierUseCase_Create(t *testing.T) {
	tests := []struct {
		name         string
		request      model.Courier
		prepare      func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller)
		expectations func(t *testing.T, id int64, err error)
	}{
		{
			name: "success: courier created",
			request: model.Courier{
				Name:          "John Doe",
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				calculator := NewMockDeliveryCalculator(ctrl)
				factory.EXPECT().
					GetDeliveryCalculator(model.TransportTypeCar).
					Return(calculator)

				repo.EXPECT().
					ExistsCourierByPhone(gomock.Any(), "+79991234567").
					Return(false, nil)

				repo.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), id)
			},
		},
		{
			name: "error: missing name",
			request: model.Courier{
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrInvalidCreate, err)
			},
		},
		{
			name: "error: missing phone",
			request: model.Courier{
				Name:          "John",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrInvalidCreate, err)
			},
		},
		{
			name: "error: missing status",
			request: model.Courier{
				Name:          "John",
				Phone:         "+79991234567",
				TransportType: "car",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrInvalidCreate, err)
			},
		},
		{
			name: "error: missing transport_type",
			request: model.Courier{
				Name:   "John",
				Phone:  "+79991234567",
				Status: "available",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrInvalidCreate, err)
			},
		},
		{
			name: "error: invalid transport type",
			request: model.Courier{
				Name:          "John Doe",
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "airplane",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				factory.EXPECT().
					GetDeliveryCalculator(model.CourierTransportType("airplane")).
					Return(nil)
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrUnknownTransportType, err)
			},
		},
		{
			name: "error: invalid phone",
			request: model.Courier{
				Name:          "John Doe",
				Phone:         "123",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				calculator := NewMockDeliveryCalculator(ctrl)
				factory.EXPECT().
					GetDeliveryCalculator(model.TransportTypeCar).
					Return(calculator)
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrInvalidPhoneNumber, err)
			},
		},
		{
			name: "error: phone already exists",
			request: model.Courier{
				Name:          "John Doe",
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				calculator := NewMockDeliveryCalculator(ctrl)
				factory.EXPECT().
					GetDeliveryCalculator(model.TransportTypeCar).
					Return(calculator)

				repo.EXPECT().
					ExistsCourierByPhone(gomock.Any(), "+79991234567").
					Return(true, nil)
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, courier.ErrPhoneNumberExists, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockcourierRepository(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
			logger, err := logger.New(logger.LogLevelInfo)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}
			uc := courier.NewCourierUseCase(mockRepo, mockFactory, logger)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockRepo, mockFactory, ctrl)
			}

			id, err := uc.CreateCourier(ctx, tc.request)

			if tc.expectations != nil {
				tc.expectations(t, id, err)
			}
		})
	}
}

func TestCourierUseCase_Update(t *testing.T) {
	nameUpdate := "Jane Doe"
	phoneUpdate := "+79991234567"
	invalidPhone := "invalid"
	transportTypeUpdate := "rocket"

	tests := []struct {
		name         string
		request      model.Courier
		prepare      func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller)
		expectations func(t *testing.T, err error)
	}{
		{
			name: "success: courier updated",
			request: model.Courier{
				ID:   1,
				Name: nameUpdate,
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				repo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error: no fields to update",
			request: model.Courier{
				ID: 1,
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, courier.ErrInvalidUpdate, err)
			},
		},
		{
			name: "error: invalid transport type",
			request: model.Courier{
				ID:            1,
				TransportType: model.CourierTransportType(transportTypeUpdate),
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				factory.EXPECT().
					GetDeliveryCalculator(model.CourierTransportType(transportTypeUpdate)).
					Return(nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, courier.ErrUnknownTransportType, err)
			},
		},
		{
			name: "error: invalid phone",
			request: model.Courier{
				ID:    1,
				Phone: invalidPhone,
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, courier.ErrInvalidPhoneNumber, err)
			},
		},
		{
			name: "error: phone already exists",
			request: model.Courier{
				ID:    1,
				Phone: phoneUpdate,
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				repo.EXPECT().
					ExistsCourierByPhone(gomock.Any(), phoneUpdate).
					Return(true, nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, courier.ErrPhoneNumberExists, err)
			},
		},
		{
			name: "error: courier not found",
			request: model.Courier{
				ID:   999,
				Name: nameUpdate,
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				repo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(courierRepo.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, courier.ErrCourierNotFound, err)
			},
		},
		{
			name: "error: repository error",
			request: model.Courier{
				ID:   1,
				Name: nameUpdate,
			},
			prepare: func(repo *MockcourierRepository, factory *MockdeliveryCalculatorFactory, ctrl *gomock.Controller) {
				repo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "database error", err.Error())
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockcourierRepository(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
			logger, err := logger.New(logger.LogLevelInfo)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}
			uc := courier.NewCourierUseCase(mockRepo, mockFactory, logger)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockRepo, mockFactory, ctrl)
			}

			err = uc.UpdateCourier(ctx, tc.request)

			if tc.expectations != nil {
				tc.expectations(t, err)
			}
		})
	}
}

func TestValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"valid phone", "+79991234567", true},
		{"invalid - no plus", "79991234567", false},
		{"invalid - too short", "+7999123456", false},
		{"invalid - too long", "+799912345678", false},
		{"invalid - letters", "+7999abc4567", false},
		{"invalid - empty", "", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := courier.ValidPhoneNumber(tc.phone)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCheckFreeCouriers(t *testing.T) {
	tests := []struct {
		name              string
		tickerInterval    time.Duration
		runDuration       time.Duration
		cancelImmediately bool
		prepare           func(repo *MockcourierRepository)
		expectations      func(t *testing.T)
	}{
		{
			name:              "success: ticker fires multiple times",
			tickerInterval:    50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(repo *MockcourierRepository) {
				repo.EXPECT().
					FreeCouriersWithInterval(gomock.Any()).
					Return(nil).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
			},
		},
		{
			name:              "repository error: continues running",
			tickerInterval:    50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(repo *MockcourierRepository) {
				repo.EXPECT().
					FreeCouriersWithInterval(gomock.Any()).
					Return(courierRepo.ErrCouriersBusy).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
			},
		},
		{
			name:              "context cancellation: stops immediately",
			tickerInterval:    50 * time.Millisecond,
			runDuration:       0,
			cancelImmediately: true,
			prepare: func(repo *MockcourierRepository) {
			},
			expectations: func(t *testing.T) {
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCourierRepo := NewMockcourierRepository(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
			logger, err := logger.New(logger.LogLevelInfo)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}
			uc := courier.NewCourierUseCase(mockCourierRepo, mockFactory, logger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tc.prepare != nil {
				tc.prepare(mockCourierRepo)
			}

			go uc.CheckFreeCouriersWithInterval(ctx, tc.tickerInterval)

			if tc.cancelImmediately {
				cancel()
			} else {
				time.Sleep(tc.runDuration)
				cancel()
			}

			time.Sleep(50 * time.Millisecond)

			if tc.expectations != nil {
				tc.expectations(t)
			}
		})
	}
}
