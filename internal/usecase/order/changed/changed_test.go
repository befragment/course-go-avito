package changed_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"courier-service/internal/model"
	"courier-service/internal/usecase/order/changed"
)

func TestOrderChangedUseCase_HandleOrderStatusChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       model.OrderStatus
		orderID      string
		prepare      func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor)
		expectations func(t *testing.T, err error)
	}{
		{
			name:    "success: completed status without gateway check",
			status:  model.OrderStatusCompleted,
			orderID: "550e8400-e29b-41d4-a716-446655440001",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				factory.EXPECT().
					Get(model.OrderStatusCompleted).
					Return(processor, true)

				processor.EXPECT().
					HandleOrderStatusChanged(gomock.Any(), model.OrderStatusCompleted, "550e8400-e29b-41d4-a716-446655440001").
					Return(nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "success: created status with gateway check",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440002",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				logger.EXPECT().
					Debugf("sending grpc request for checking status for order %s", "550e8400-e29b-41d4-a716-446655440002")

				gateway.EXPECT().
					GetOrderById(gomock.Any(), "550e8400-e29b-41d4-a716-446655440002").
					Return(model.Order{
						ID:     "550e8400-e29b-41d4-a716-446655440002",
						Status: model.OrderStatusCreated,
					}, nil)

				factory.EXPECT().
					Get(model.OrderStatusCreated).
					Return(processor, true)

				processor.EXPECT().
					HandleOrderStatusChanged(gomock.Any(), model.OrderStatusCreated, "550e8400-e29b-41d4-a716-446655440002").
					Return(nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "success: cancelled status with gateway check",
			status:  model.OrderStatusCancelled,
			orderID: "550e8400-e29b-41d4-a716-446655440003",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				logger.EXPECT().
					Debugf("sending grpc request for checking status for order %s", "550e8400-e29b-41d4-a716-446655440003")

				gateway.EXPECT().
					GetOrderById(gomock.Any(), "550e8400-e29b-41d4-a716-446655440003").
					Return(model.Order{
						ID:     "550e8400-e29b-41d4-a716-446655440003",
						Status: model.OrderStatusCancelled,
					}, nil)

				factory.EXPECT().
					Get(model.OrderStatusCancelled).
					Return(processor, true)

				processor.EXPECT().
					HandleOrderStatusChanged(gomock.Any(), model.OrderStatusCancelled, "550e8400-e29b-41d4-a716-446655440003").
					Return(nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "success: no processor for status",
			status:  model.OrderStatus("in_progress"),
			orderID: "550e8400-e29b-41d4-a716-446655440004",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				logger.EXPECT().
					Debugf("sending grpc request for checking status for order %s", "550e8400-e29b-41d4-a716-446655440004")

				gateway.EXPECT().
					GetOrderById(gomock.Any(), "550e8400-e29b-41d4-a716-446655440004").
					Return(model.Order{
						ID:     "550e8400-e29b-41d4-a716-446655440004",
						Status: model.OrderStatus("in_progress"),
					}, nil)

				factory.EXPECT().
					Get(model.OrderStatus("in_progress")).
					Return(nil, false)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "error: gateway returns error",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440005",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				logger.EXPECT().
					Debugf("sending grpc request for checking status for order %s", "550e8400-e29b-41d4-a716-446655440005")

				gateway.EXPECT().
					GetOrderById(gomock.Any(), "550e8400-e29b-41d4-a716-446655440005").
					Return(model.Order{}, errors.New("gateway error"))
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "gateway error")
			},
		},
		{
			name:    "error: status mismatch",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440006",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				logger.EXPECT().
					Debugf("sending grpc request for checking status for order %s", "550e8400-e29b-41d4-a716-446655440006")

				gateway.EXPECT().
					GetOrderById(gomock.Any(), "550e8400-e29b-41d4-a716-446655440006").
					Return(model.Order{
						ID:     "550e8400-e29b-41d4-a716-446655440006",
						Status: model.OrderStatusCancelled, // Different status
					}, nil)

				logger.EXPECT().
					Warnf("order status mismatch: expected %s, got %s for order %s",
						model.OrderStatusCreated,
						model.OrderStatusCancelled,
						"550e8400-e29b-41d4-a716-446655440006")
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, changed.ErrOrderStatusMismatch, err)
			},
		},
		{
			name:    "error: processor returns error",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440007",
			prepare: func(factory *MockorderChangedFactory, gateway *MockorderGateway, logger *Mocklogger, processor *MockProcessor) {
				logger.EXPECT().
					Debugf("sending grpc request for checking status for order %s", "550e8400-e29b-41d4-a716-446655440007")

				gateway.EXPECT().
					GetOrderById(gomock.Any(), "550e8400-e29b-41d4-a716-446655440007").
					Return(model.Order{
						ID:     "550e8400-e29b-41d4-a716-446655440007",
						Status: model.OrderStatusCreated,
					}, nil)

				factory.EXPECT().
					Get(model.OrderStatusCreated).
					Return(processor, true)

				processor.EXPECT().
					HandleOrderStatusChanged(gomock.Any(), model.OrderStatusCreated, "550e8400-e29b-41d4-a716-446655440007").
					Return(errors.New("processor error"))
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "processor error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := NewMockorderChangedFactory(ctrl)
			mockGateway := NewMockorderGateway(ctrl)
			mockLogger := NewMocklogger(ctrl)
			mockProcessor := NewMockProcessor(ctrl)

			uc := changed.NewOrderChangedUseCase(mockFactory, mockGateway, mockLogger)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockFactory, mockGateway, mockLogger, mockProcessor)
			}

			err := uc.HandleOrderStatusChanged(ctx, tc.status, tc.orderID)

			if tc.expectations != nil {
				tc.expectations(t, err)
			}
		})
	}
}
