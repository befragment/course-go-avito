package processor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"courier-service/internal/model"
	"courier-service/internal/usecase/delivery/assign"
	complete "courier-service/internal/usecase/delivery/complete"
	changed "courier-service/internal/usecase/order/changed"
	"courier-service/internal/usecase/order/changed/processor"
)

func TestCreatedProcessor_HandleOrderStatusChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       model.OrderStatus
		orderID      string
		prepare      func(assignUC *MockassignUseCase)
		expectations func(t *testing.T, err error)
	}{
		{
			name:    "success: order assigned to courier",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440001",
			prepare: func(assignUC *MockassignUseCase) {
				assignUC.EXPECT().
					Assign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440001").
					Return(assign.DeliveryAssignResponse{
						CourierID:     1,
						OrderID:       "550e8400-e29b-41d4-a716-446655440001",
						TransportType: "car",
					}, nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "error: assign use case returns error",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440002",
			prepare: func(assignUC *MockassignUseCase) {
				assignUC.EXPECT().
					Assign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440002").
					Return(assign.DeliveryAssignResponse{}, errors.New("no available couriers"))
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "no available couriers")
			},
		},
		{
			name:    "error: assign use case returns ErrCouriersBusy",
			status:  model.OrderStatusCreated,
			orderID: "550e8400-e29b-41d4-a716-446655440003",
			prepare: func(assignUC *MockassignUseCase) {
				assignUC.EXPECT().
					Assign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440003").
					Return(assign.DeliveryAssignResponse{}, assign.ErrCouriersBusy)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, assign.ErrCouriersBusy, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAssignUC := NewMockassignUseCase(ctrl)
			proc := processor.NewCreatedProcessor(mockAssignUC)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockAssignUC)
			}

			err := proc.HandleOrderStatusChanged(ctx, tc.status, tc.orderID)

			if tc.expectations != nil {
				tc.expectations(t, err)
			}
		})
	}
}

func TestCancelledProcessor_HandleOrderStatusChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       model.OrderStatus
		orderID      string
		prepare      func(unassignUC *MockunassignUseCase)
		expectations func(t *testing.T, err error)
	}{
		{
			name:    "success: order unassigned from courier",
			status:  model.OrderStatusCancelled,
			orderID: "550e8400-e29b-41d4-a716-446655440001",
			prepare: func(unassignUC *MockunassignUseCase) {
				unassignUC.EXPECT().
					Unassign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440001").
					Return(int64(1), nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "error: unassign use case returns error",
			status:  model.OrderStatusCancelled,
			orderID: "550e8400-e29b-41d4-a716-446655440002",
			prepare: func(unassignUC *MockunassignUseCase) {
				unassignUC.EXPECT().
					Unassign(gomock.Any(), "550e8400-e29b-41d4-a716-446655440002").
					Return(int64(0), errors.New("delivery not found"))
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "delivery not found")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUnassignUC := NewMockunassignUseCase(ctrl)
			proc := processor.NewCancelledProcessor(mockUnassignUC)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockUnassignUC)
			}

			err := proc.HandleOrderStatusChanged(ctx, tc.status, tc.orderID)

			if tc.expectations != nil {
				tc.expectations(t, err)
			}
		})
	}
}

func TestCompletedProcessor_HandleOrderStatusChanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       model.OrderStatus
		orderID      string
		prepare      func(completeUC *MockcompleteUseCase)
		expectations func(t *testing.T, err error)
	}{
		{
			name:    "success: delivery completed",
			status:  model.OrderStatusCompleted,
			orderID: "550e8400-e29b-41d4-a716-446655440001",
			prepare: func(completeUC *MockcompleteUseCase) {
				completeUC.EXPECT().
					Complete(gomock.Any(), "550e8400-e29b-41d4-a716-446655440001").
					Return(nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:    "error: order not found",
			status:  model.OrderStatusCompleted,
			orderID: "550e8400-e29b-41d4-a716-446655440002",
			prepare: func(completeUC *MockcompleteUseCase) {
				completeUC.EXPECT().
					Complete(gomock.Any(), "550e8400-e29b-41d4-a716-446655440002").
					Return(complete.ErrOrderNotFound)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, changed.ErrOrderNotFound, err)
			},
		},
		{
			name:    "error: database error",
			status:  model.OrderStatusCompleted,
			orderID: "550e8400-e29b-41d4-a716-446655440003",
			prepare: func(completeUC *MockcompleteUseCase) {
				completeUC.EXPECT().
					Complete(gomock.Any(), "550e8400-e29b-41d4-a716-446655440003").
					Return(errors.New("database connection lost"))
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database error")
				assert.Contains(t, err.Error(), "database connection lost")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompleteUC := NewMockcompleteUseCase(ctrl)
			proc := processor.NewCompletedProcessor(mockCompleteUC)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockCompleteUC)
			}

			err := proc.HandleOrderStatusChanged(ctx, tc.status, tc.orderID)

			if tc.expectations != nil {
				tc.expectations(t, err)
			}
		})
	}
}
