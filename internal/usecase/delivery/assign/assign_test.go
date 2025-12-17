package assign_test

import (
	"context"
	"testing"
	"time"

	"courier-service/internal/model"
	courierstorage "courier-service/internal/repository/courier"
	deliverystorage "courier-service/internal/repository/delivery"
	"courier-service/internal/usecase/delivery/assign"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAssignDelivery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		orderID string
		prepare func(
			courierRepository *MockcourierRepository,
			deliveryRepository *MockdeliveryRepository,
			txRunner *MocktxRunner,
			factory *MockdeliveryCalculatorFactory,
			ctrl *gomock.Controller,
		)
		expectations func(t *testing.T, resp assign.DeliveryAssignResponse, err error)
	}{
		{
			name:    "success: delivery assigned",
			orderID: "550e8400-e29b-41d4-a716-446655440001",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
				factory *MockdeliveryCalculatorFactory,
				ctrl *gomock.Controller,
			) {
				now := time.Now()

				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				calculator := NewMockDeliveryCalculator(ctrl)
				factory.EXPECT().
					GetDeliveryCalculator(model.TransportTypeCar).
					Return(calculator)
				calculator.EXPECT().
					CalculateDeadline().
					Return(now.Add(5 * time.Minute))

				courierRepository.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        model.CourierStatusAvailable,
						TransportType: "car",
					}, nil)

				deliveryRepository.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(model.Delivery{
						ID:         1,
						CourierID:  1,
						OrderID:    "550e8400-e29b-41d4-a716-446655440001",
						AssignedAt: now,
						Deadline:   now.Add(5 * time.Minute),
					}, nil)

				courierRepository.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, c model.Courier) error {
						assert.Equal(t, model.CourierStatusBusy, c.Status)
						return nil
					})
			},
			expectations: func(t *testing.T, resp assign.DeliveryAssignResponse, err error) {
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
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
				factory *MockdeliveryCalculatorFactory,
				ctrl *gomock.Controller,
			) {
				// No mock expectations - validation happens before any repo calls
			},
			expectations: func(t *testing.T, resp assign.DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, assign.ErrNoOrderID, err)
				assert.Equal(t, assign.DeliveryAssignResponse{}, resp)
			},
		},
		{
			name:    "error: all couriers busy",
			orderID: "550e8400-e29b-41d4-a716-446655440002",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
				factory *MockdeliveryCalculatorFactory,
				ctrl *gomock.Controller,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				courierRepository.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{}, courierstorage.ErrCouriersBusy)
			},
			expectations: func(t *testing.T, resp assign.DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, assign.ErrCouriersBusy, err)
				assert.Equal(t, assign.DeliveryAssignResponse{}, resp)
			},
		},
		{
			name:    "error: order ID already exists",
			orderID: "550e8400-e29b-41d4-a716-446655440003",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
				factory *MockdeliveryCalculatorFactory,
				ctrl *gomock.Controller,
			) {
				now := time.Now()

				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				calculator := NewMockDeliveryCalculator(ctrl)

				factory.EXPECT().
					GetDeliveryCalculator(model.TransportTypeCar).
					Return(calculator)
				calculator.EXPECT().
					CalculateDeadline().
					Return(now.Add(5 * time.Minute))

				courierRepository.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        model.CourierStatusAvailable,
						TransportType: "car",
					}, nil)

				deliveryRepository.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(model.Delivery{}, deliverystorage.ErrOrderIDExists)
			},
			expectations: func(t *testing.T, resp assign.DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, assign.ErrOrderIDExists, err)
				assert.Equal(t, assign.DeliveryAssignResponse{}, resp)
			},
		},
		{
			name:    "error: failed to update courier status",
			orderID: "550e8400-e29b-41d4-a716-446655440004",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
				factory *MockdeliveryCalculatorFactory,
				ctrl *gomock.Controller,
			) {
				now := time.Now()

				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				calculator := NewMockDeliveryCalculator(ctrl)
				factory.EXPECT().
					GetDeliveryCalculator(model.TransportTypeCar).
					Return(calculator)
				calculator.EXPECT().
					CalculateDeadline().
					Return(now.Add(5 * time.Minute))

				courierRepository.EXPECT().
					FindAvailableCourier(gomock.Any()).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)

				deliveryRepository.EXPECT().
					CreateDelivery(gomock.Any(), gomock.Any()).
					Return(model.Delivery{
						ID:         1,
						CourierID:  1,
						OrderID:    "550e8400-e29b-41d4-a716-446655440004",
						AssignedAt: now,
						Deadline:   now.Add(5 * time.Minute),
					}, nil)

				courierRepository.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(courierstorage.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, resp assign.DeliveryAssignResponse, err error) {
				assert.Error(t, err)
				assert.Equal(t, courierstorage.ErrCourierNotFound, err)
				assert.Equal(t, assign.DeliveryAssignResponse{}, resp)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCourierRepo := NewMockcourierRepository(ctrl)
			mockDeliveryRepo := NewMockdeliveryRepository(ctrl)
			mockTxRunner := NewMocktxRunner(ctrl)
			mockFactory := NewMockdeliveryCalculatorFactory(ctrl)

			uc := assign.NewAssignDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner, mockFactory)

			ctx := context.Background()
			req := assign.DeliveryAssignRequest{
				OrderID: tc.orderID,
			}

			if tc.prepare != nil {
				tc.prepare(mockCourierRepo, mockDeliveryRepo, mockTxRunner, mockFactory, ctrl)
			}

			result, err := uc.Assign(ctx, req)

			if tc.expectations != nil {
				tc.expectations(t, result, err)
			}
		})
	}
}
