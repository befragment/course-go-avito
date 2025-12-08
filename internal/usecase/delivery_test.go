package usecase

import (
	"context"
	"testing"
	"time"

	"courier-service/internal/model"
	"courier-service/internal/repository"
	"courier-service/internal/usecase/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestDeliveryUseCase_AssignDelivery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		orderID string
		prepare func(
			courierRepo *mocks.MockсourierRepository,
			deliveryRepo *mocks.MockdeliveryRepository,
			txRunner *mocks.MocktxRunner,
		)
		expectations func(t *testing.T, resp DeliveryAssignResponse, err error)
	}{
		{
			name:    "success: delivery assigned",
			orderID: "550e8400-e29b-41d4-a716-446655440001",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				courierRepo.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)

				now := time.Now()
				deliveryRepo.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(model.Delivery{
						ID:         1,
						CourierID:  1,
						OrderID:    "550e8400-e29b-41d4-a716-446655440001",
						AssignedAt: now,
						Deadline:   now.Add(5 * time.Minute),
					}, nil)

				courierRepo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, c model.Courier) error {
						assert.Equal(t, "busy", c.Status)
						return nil
					})
			},
			expectations: func(t *testing.T, resp DeliveryAssignResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), resp.CourierID)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", resp.OrderID)
				assert.Equal(t, "car", resp.TransportType)
			},
		},
		{
			name:    "error: no order ID",
			orderID: "",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				// No mock expectations - validation happens before any repo calls
			},
			expectations: func(t *testing.T, resp DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrNoOrderID, err)
				assert.Equal(t, DeliveryAssignResponse{}, resp)
			},
		},
		{
			name:    "error: all couriers busy",
			orderID: "550e8400-e29b-41d4-a716-446655440002",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				courierRepo.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{}, repository.ErrCouriersBusy)
			},
			expectations: func(t *testing.T, resp DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrCouriersBusy, err)
				assert.Equal(t, DeliveryAssignResponse{}, resp)
			},
		},
		{
			name:    "error: order ID already exists",
			orderID: "550e8400-e29b-41d4-a716-446655440003",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				courierRepo.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)

				deliveryRepo.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(model.Delivery{}, repository.ErrOrderIDExists)
			},
			expectations: func(t *testing.T, resp DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrOrderIDExists, err)
				assert.Equal(t, DeliveryAssignResponse{}, resp)
			},
		},
		{
			name:    "error: failed to update courier status",
			orderID: "550e8400-e29b-41d4-a716-446655440004",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				courierRepo.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)

				now := time.Now()
				deliveryRepo.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(model.Delivery{
						ID:         1,
						CourierID:  1,
						OrderID:    "550e8400-e29b-41d4-a716-446655440004",
						AssignedAt: now,
						Deadline:   now.Add(5 * time.Minute),
					}, nil)

				courierRepo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(repository.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, resp DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, repository.ErrCourierNotFound, err)
				assert.Equal(t, DeliveryAssignResponse{}, resp)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
			mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
			mockTxRunner := mocks.NewMocktxRunner(ctrl)

			uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

			ctx := context.Background()
			req := DeliveryAssignRequest{
				OrderID: tc.orderID,
			}

			if tc.prepare != nil {
				tc.prepare(mockCourierRepo, mockDeliveryRepo, mockTxRunner)
			}

			result, err := uc.AssignDelivery(ctx, req)

			if tc.expectations != nil {
				tc.expectations(t, result, err)
			}
		})
	}
}

func TestDeliveryUseCase_UnassignDelivery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		orderID string
		prepare func(
			courierRepo *mocks.MockсourierRepository,
			deliveryRepo *mocks.MockdeliveryRepository,
			txRunner *mocks.MocktxRunner,
		)
		expectations func(t *testing.T, resp DeliveryUnassignResponse, err error)
	}{
		{
			name:    "success: delivery unassigned",
			orderID: "550e8400-e29b-41d4-a716-446655440005",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepo.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440005").
					Return(model.Delivery{
						ID:        1,
						CourierID: 1,
						OrderID:   "550e8400-e29b-41d4-a716-446655440005",
					}, nil)

				deliveryRepo.EXPECT().
					DeleteDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440005").
					Return(nil)

				courierRepo.EXPECT().
					GetCourierById(gomock.Any(), int64(1)).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        "busy",
						TransportType: "car",
					}, nil)

				courierRepo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, c model.Courier) error {
						assert.Equal(t, "available", c.Status)
						return nil
					})
			},
			expectations: func(t *testing.T, resp DeliveryUnassignResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440005", resp.OrderID)
				assert.Equal(t, int64(1), resp.CourierID)
				assert.Equal(t, "unassigned", resp.Status)
			},
		},
		{
			name:    "error: no order ID",
			orderID: "",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				// No mock expectations - validation happens before any repo calls
			},
			expectations: func(t *testing.T, resp DeliveryUnassignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrNoOrderID, err)
				assert.Equal(t, DeliveryUnassignResponse{}, resp)
			},
		},
		{
			name:    "error: order not found",
			orderID: "550e8400-e29b-41d4-a716-446655440006",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepo.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440006").
					Return(model.Delivery{}, repository.ErrOrderIDNotFound)
			},
			expectations: func(t *testing.T, resp DeliveryUnassignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrOrderIDNotFound, err)
				assert.Equal(t, DeliveryUnassignResponse{}, resp)
			},
		},
		{
			name:    "error: failed to delete delivery",
			orderID: "550e8400-e29b-41d4-a716-446655440007",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepo.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440007").
					Return(model.Delivery{
						ID:        1,
						CourierID: 1,
						OrderID:   "550e8400-e29b-41d4-a716-446655440007",
					}, nil)

				deliveryRepo.EXPECT().
					DeleteDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440007").
					Return(repository.ErrOrderIDNotFound)
			},
			expectations: func(t *testing.T, resp DeliveryUnassignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrOrderIDNotFound, err)
				assert.Equal(t, DeliveryUnassignResponse{}, resp)
			},
		},
		{
			name:    "error: courier not found",
			orderID: "550e8400-e29b-41d4-a716-446655440008",
			prepare: func(
				courierRepo *mocks.MockсourierRepository,
				deliveryRepo *mocks.MockdeliveryRepository,
				txRunner *mocks.MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepo.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440008").
					Return(model.Delivery{
						ID:        1,
						CourierID: 999,
						OrderID:   "550e8400-e29b-41d4-a716-446655440008",
					}, nil)

				deliveryRepo.EXPECT().
					DeleteDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440008").
					Return(nil)

				courierRepo.EXPECT().
					GetCourierById(gomock.Any(), int64(999)).
					Return(model.Courier{}, repository.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, resp DeliveryUnassignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, repository.ErrCourierNotFound, err)
				assert.Equal(t, DeliveryUnassignResponse{}, resp)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
			mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
			mockTxRunner := mocks.NewMocktxRunner(ctrl)

			uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

			ctx := context.Background()
			req := DeliveryUnassignRequest{
				OrderID: tc.orderID,
			}

			if tc.prepare != nil {
				tc.prepare(mockCourierRepo, mockDeliveryRepo, mockTxRunner)
			}

			resp, err := uc.UnassignDelivery(ctx, req)

			if tc.expectations != nil {
				tc.expectations(t, resp, err)
			}
		})
	}
}

func TestTransportTypeTime(t *testing.T) {
	testCases := []struct {
		name          string
		transportType string
		expected      time.Duration
		expectError   bool
	}{
		{"car", "car", 5 * time.Minute, false},
		{"scooter", "scooter", 15 * time.Minute, false},
		{"on_foot", "on_foot", 30 * time.Minute, false},
		{"invalid", "rocket", 0, true},
		{"empty", "", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			duration, err := transportTypeTime(tc.transportType)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrUnknownTransportType, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, duration)
		})
	}
}

func TestCheckFreeCouriers(t *testing.T) {
	tests := []struct {
		name              string
		tickerInterval    time.Duration
		runDuration       time.Duration
		cancelImmediately bool
		prepare           func(repo *mocks.MockсourierRepository)
		expectations      func(t *testing.T)
	}{
		{
			name:              "success: ticker fires multiple times",
			tickerInterval:    50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					FreeCouriersWithInterval(gomock.Any()).
					Return(nil).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
				// goleak will verify no goroutines leaked
			},
		},
		{
			name:              "repository error: continues running",
			tickerInterval:    50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					FreeCouriersWithInterval(gomock.Any()).
					Return(repository.ErrCouriersBusy).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
				// goleak will verify no goroutines leaked
			},
		},
		{
			name:              "context cancellation: stops immediately",
			tickerInterval:    50 * time.Millisecond,
			runDuration:       0,
			cancelImmediately: true,
			prepare: func(repo *mocks.MockсourierRepository) {
				// No expectations - context should be cancelled before ticker fires
			},
			expectations: func(t *testing.T) {
				// goleak will verify no goroutines leaked
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
			uc := NewCourierUseCase(mockCourierRepo)

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
