package unassign_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"courier-service/internal/model"
	courierstorage "courier-service/internal/repository/courier"
	deliverystorage "courier-service/internal/repository/delivery"
	"courier-service/internal/usecase/delivery/unassign"
)

func TestUnassignDelivery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		orderID string
		prepare func(
			courierRepository *MockcourierRepository,
			deliveryRepository *MockdeliveryRepository,
			txRunner *MocktxRunner,
		)
		expectations func(t *testing.T, resp int64, err error)
	}{
		{
			name:    "success: delivery unassigned",
			orderID: "550e8400-e29b-41d4-a716-446655440005",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepository.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440005").
					Return(model.Delivery{
						ID:        1,
						CourierID: 1,
						OrderID:   "550e8400-e29b-41d4-a716-446655440005",
					}, nil)

				deliveryRepository.EXPECT().
					DeleteDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440005").
					Return(nil)

				courierRepository.EXPECT().
					GetCourierById(gomock.Any(), int64(1)).
					Return(model.Courier{
						ID:            1,
						Name:          "John",
						Phone:         "+79991234567",
						Status:        model.CourierStatusBusy,
						TransportType: "car",
					}, nil)

				courierRepository.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, c model.Courier) error {
						assert.Equal(t, model.CourierStatusAvailable, c.Status)
						return nil
					})
			},
			expectations: func(t *testing.T, resp int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), resp)
			},
		},
		{
			name:    "error: no order ID",
			orderID: "",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
			) {
				// No mock expectations - validation happens before any repo calls
			},
			expectations: func(t *testing.T, resp int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, unassign.ErrNoOrderID, err)
				assert.Equal(t, int64(0), resp)
			},
		},
		{
			name:    "error: order not found",
			orderID: "550e8400-e29b-41d4-a716-446655440006",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepository.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440006").
					Return(model.Delivery{}, deliverystorage.ErrOrderIDNotFound)
			},
			expectations: func(t *testing.T, resp int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, unassign.ErrOrderIDNotFound, err)
				assert.Equal(t, int64(0), resp)
			},
		},
		{
			name:    "error: failed to delete delivery",
			orderID: "550e8400-e29b-41d4-a716-446655440007",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepository.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440007").
					Return(model.Delivery{
						ID:        1,
						CourierID: 1,
						OrderID:   "550e8400-e29b-41d4-a716-446655440007",
					}, nil)

				deliveryRepository.EXPECT().
					DeleteDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440007").
					Return(deliverystorage.ErrOrderIDNotFound)
			},
			expectations: func(t *testing.T, resp int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, unassign.ErrOrderIDNotFound, err)
				assert.Equal(t, int64(0), resp)
			},
		},
		{
			name:    "error: courier not found",
			orderID: "550e8400-e29b-41d4-a716-446655440008",
			prepare: func(
				courierRepository *MockcourierRepository,
				deliveryRepository *MockdeliveryRepository,
				txRunner *MocktxRunner,
			) {
				txRunner.EXPECT().
					Run(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				deliveryRepository.EXPECT().
					CouriersDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440008").
					Return(model.Delivery{
						ID:        1,
						CourierID: 999,
						OrderID:   "550e8400-e29b-41d4-a716-446655440008",
					}, nil)

				deliveryRepository.EXPECT().
					DeleteDelivery(gomock.Any(), "550e8400-e29b-41d4-a716-446655440008").
					Return(nil)

				courierRepository.EXPECT().
					GetCourierById(gomock.Any(), int64(999)).
					Return(model.Courier{}, courierstorage.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, resp int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, courierstorage.ErrCourierNotFound, err)
				assert.Equal(t, int64(0), resp)
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

			if tc.prepare != nil {
				tc.prepare(mockCourierRepo, mockDeliveryRepo, mockTxRunner)
			}

			uc := unassign.NewUnassignDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

			ctx := context.Background()

			resp, err := uc.Unassign(ctx, tc.orderID)

			if tc.expectations != nil {
				tc.expectations(t, resp, err)
			}
		})
	}
}
