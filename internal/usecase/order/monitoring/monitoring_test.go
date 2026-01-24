package ordermonitoring_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"courier-service/internal/model"
	"courier-service/internal/usecase/delivery/assign"
	ordermonitoring "courier-service/internal/usecase/order/monitoring"
)

func TestOrderMonitoringUseCase_MonitorOrders(t *testing.T) {
	tests := []struct {
		name              string
		interval          time.Duration
		runDuration       time.Duration
		cancelImmediately bool
		prepare           func(gateway *MockorderGateway, assignUC *MockassignUseCase)
		expectations      func(t *testing.T)
	}{
		{
			name:              "success: processes orders from gateway",
			interval:          50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(gateway *MockorderGateway, assignUC *MockassignUseCase) {
				gateway.EXPECT().
					GetOrders(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, from time.Time) ([]model.Order, error) {
						return []model.Order{
							{
								ID:     "550e8400-e29b-41d4-a716-446655440001",
								Status: model.OrderStatusCreated,
							},
							{
								ID:     "550e8400-e29b-41d4-a716-446655440002",
								Status: model.OrderStatusCreated,
							},
						}, nil
					}).
					MinTimes(2)

				assignUC.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, orderID string) (assign.DeliveryAssignResponse, error) {
						return assign.DeliveryAssignResponse{
							CourierID:     1,
							OrderID:       orderID,
							TransportType: "car",
						}, nil
					}).
					MinTimes(4)
			},
			expectations: func(t *testing.T) {
				// Success case
			},
		},
		{
			name:              "success: handles empty order list",
			interval:          50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(gateway *MockorderGateway, assignUC *MockassignUseCase) {
				gateway.EXPECT().
					GetOrders(gomock.Any(), gomock.Any()).
					Return([]model.Order{}, nil).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
				// Success case with empty orders
			},
		},
		{
			name:              "error: gateway returns error, continues monitoring",
			interval:          50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(gateway *MockorderGateway, assignUC *MockassignUseCase) {
				gateway.EXPECT().
					GetOrders(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("gateway connection failed")).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
				// Error should be logged but monitoring continues
			},
		},
		{
			name:              "error: assign fails for some orders, continues with others",
			interval:          50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(gateway *MockorderGateway, assignUC *MockassignUseCase) {
				gateway.EXPECT().
					GetOrders(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, from time.Time) ([]model.Order, error) {
						return []model.Order{
							{
								ID:     "550e8400-e29b-41d4-a716-446655440001",
								Status: model.OrderStatusCreated,
							},
							{
								ID:     "550e8400-e29b-41d4-a716-446655440002",
								Status: model.OrderStatusCreated,
							},
						}, nil
					}).
					MinTimes(2)

				assignUC.EXPECT().
					Assign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440001").
					Return(assign.DeliveryAssignResponse{}, errors.New("no available couriers")).
					MinTimes(2)

				assignUC.EXPECT().
					Assign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440002").
					Return(assign.DeliveryAssignResponse{
						CourierID:     1,
						OrderID:       "550e8400-e29b-41d4-a716-446655440002",
						TransportType: "car",
					}, nil).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
				// Some assignments fail, but monitoring continues
			},
		},
		{
			name:              "success: context cancellation stops monitoring immediately",
			interval:          50 * time.Millisecond,
			runDuration:       0,
			cancelImmediately: true,
			prepare: func(gateway *MockorderGateway, assignUC *MockassignUseCase) {
				// No expectations - context should be cancelled before first tick
			},
			expectations: func(t *testing.T) {
				// Context cancellation should stop monitoring
			},
		},
		{
			name:              "success: monitors single order successfully",
			interval:          50 * time.Millisecond,
			runDuration:       150 * time.Millisecond,
			cancelImmediately: false,
			prepare: func(gateway *MockorderGateway, assignUC *MockassignUseCase) {
				gateway.EXPECT().
					GetOrders(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, from time.Time) ([]model.Order, error) {
						return []model.Order{
							{
								ID:     "550e8400-e29b-41d4-a716-446655440001",
								Status: model.OrderStatusCreated,
							},
						}, nil
					}).
					MinTimes(2)

				assignUC.EXPECT().
					Assign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440001").
					Return(assign.DeliveryAssignResponse{
						CourierID:     1,
						OrderID:       "550e8400-e29b-41d4-a716-446655440001",
						TransportType: "car",
					}, nil).
					MinTimes(2)
			},
			expectations: func(t *testing.T) {
				// Success case
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGateway := NewMockorderGateway(ctrl)
			mockCourierRepo := NewMockcourierRepository(ctrl)
			mockDeliveryRepo := NewMockdeliveryRepository(ctrl)
			mockTxRunner := NewMocktxRunner(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
			mockAssignUC := NewMockassignUseCase(ctrl)

			uc := ordermonitoring.NewOrderMonitoringUseCase(
				mockGateway,
				mockCourierRepo,
				mockDeliveryRepo,
				mockTxRunner,
				mockFactory,
				mockAssignUC,
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tc.prepare != nil {
				tc.prepare(mockGateway, mockAssignUC)
			}

			go uc.MonitorOrders(ctx, tc.interval)

			if tc.cancelImmediately {
				cancel()
			} else {
				time.Sleep(tc.runDuration)
				cancel()
			}

			// Give some time for goroutine to finish
			time.Sleep(100 * time.Millisecond)

			if tc.expectations != nil {
				tc.expectations(t)
			}
		})
	}
}

func TestOrderMonitoringUseCase_MonitorOrders_TimeWindow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := NewMockorderGateway(ctrl)
	mockCourierRepo := NewMockcourierRepository(ctrl)
	mockDeliveryRepo := NewMockdeliveryRepository(ctrl)
	mockTxRunner := NewMocktxRunner(ctrl)
	mockFactory := NewMockdeliveryCalculatorFactory(ctrl)
	mockAssignUC := NewMockassignUseCase(ctrl)

	uc := ordermonitoring.NewOrderMonitoringUseCase(
		mockGateway,
		mockCourierRepo,
		mockDeliveryRepo,
		mockTxRunner,
		mockFactory,
		mockAssignUC,
	)

	interval := 100 * time.Millisecond

	var capturedTime time.Time
	mockGateway.EXPECT().
		GetOrders(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, from time.Time) ([]model.Order, error) {
			capturedTime = from
			return []model.Order{}, nil
		}).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()
	go uc.MonitorOrders(ctx, interval)

	// Wait for first tick
	time.Sleep(150 * time.Millisecond)
	cancel()

	// Give time for goroutine to finish
	time.Sleep(50 * time.Millisecond)

	// Verify that the time window is approximately correct
	expectedTime := startTime.Add(-interval)
	timeDiff := capturedTime.Sub(expectedTime)
	assert.True(t, timeDiff < 150*time.Millisecond && timeDiff > -150*time.Millisecond,
		"Time window should be approximately %v ago, but was %v", interval, timeDiff)
}
