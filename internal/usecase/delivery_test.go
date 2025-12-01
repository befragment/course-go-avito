package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"courier-service/internal/model"
	"courier-service/internal/repository"
	"courier-service/internal/usecase/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestDeliveryUseCase_AssignDelivery_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryAssignRequest{
		OrderID: "ORDER123",
	}

	courierDB := &model.CourierDB{
		ID:            1,
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	now := time.Now()
	delivery := &model.Delivery{
		ID:         1,
		CourierID:  1,
		OrderID:    "ORDER123",
		AssignedAt: now,
		Deadline:   now.Add(20 * time.Second),
	}

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockCourierRepo.EXPECT().
		FindAvailableCourier(ctx).
		Return(courierDB, nil)

	mockDeliveryRepo.EXPECT().
		CreateDelivery(ctx, gomock.Any()).
		Return(delivery, nil)

	mockCourierRepo.EXPECT().
		UpdateCourier(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, c *model.CourierDB) error {
			assert.Equal(t, "busy", c.Status)
			return nil
		})

	resp, err := uc.AssignDelivery(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.CourierID)
	assert.Equal(t, "ORDER123", resp.OrderID)
	assert.Equal(t, "car", resp.TransportType)
}

func TestDeliveryUseCase_AssignDelivery_NoOrderID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryAssignRequest{
		OrderID: "",
	}

	resp, err := uc.AssignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrNoOrderID, err)
	assert.Equal(t, model.DeliveryAssignResponse{}, resp)
}

func TestDeliveryUseCase_AssignDelivery_CouriersBusy(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryAssignRequest{
		OrderID: "ORDER123",
	}

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockCourierRepo.EXPECT().
		FindAvailableCourier(ctx).
		Return(nil, repository.ErrCouriersBusy)

	resp, err := uc.AssignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrCouriersBusy, err)
	assert.Equal(t, model.DeliveryAssignResponse{}, resp)
}

func TestDeliveryUseCase_AssignDelivery_OrderExists(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryAssignRequest{
		OrderID: "ORDER123",
	}

	courierDB := &model.CourierDB{
		ID:            1,
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockCourierRepo.EXPECT().
		FindAvailableCourier(ctx).
		Return(courierDB, nil)

	mockDeliveryRepo.EXPECT().
		CreateDelivery(ctx, gomock.Any()).
		Return(nil, repository.ErrOrderIDExists)

	resp, err := uc.AssignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderIDExists, err)
	assert.Equal(t, model.DeliveryAssignResponse{}, resp)
}

func TestDeliveryUseCase_AssignDelivery_UpdateCourierError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryAssignRequest{
		OrderID: "ORDER123",
	}

	courierDB := &model.CourierDB{
		ID:            1,
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	now := time.Now()
	delivery := &model.Delivery{
		ID:         1,
		CourierID:  1,
		OrderID:    "ORDER123",
		AssignedAt: now,
		Deadline:   now.Add(20 * time.Second),
	}

	expectedErr := errors.New("update error")

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockCourierRepo.EXPECT().
		FindAvailableCourier(ctx).
		Return(courierDB, nil)

	mockDeliveryRepo.EXPECT().
		CreateDelivery(ctx, gomock.Any()).
		Return(delivery, nil)

	mockCourierRepo.EXPECT().
		UpdateCourier(ctx, gomock.Any()).
		Return(expectedErr)

	resp, err := uc.AssignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, model.DeliveryAssignResponse{}, resp)
}

func TestDeliveryUseCase_UnassignDelivery_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryUnassignRequest{
		OrderID: "ORDER123",
	}

	deliveryDB := &model.DeliveryDB{
		ID:        1,
		CourierID: 1,
		OrderID:   "ORDER123",
	}

	courierDB := &model.CourierDB{
		ID:            1,
		Name:          "John",
		Phone:         "+79991234567",
		Status:        "busy",
		TransportType: "car",
	}

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockDeliveryRepo.EXPECT().
		CouriersDelivery(ctx, "ORDER123").
		Return(deliveryDB, nil)

	mockDeliveryRepo.EXPECT().
		DeleteDelivery(ctx, "ORDER123").
		Return(nil)

	mockCourierRepo.EXPECT().
		GetCourierById(ctx, int64(1)).
		Return(courierDB, nil)

	mockCourierRepo.EXPECT().
		UpdateCourier(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, c *model.CourierDB) error {
			assert.Equal(t, "available", c.Status)
			return nil
		})

	resp, err := uc.UnassignDelivery(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "ORDER123", resp.OrderID)
	assert.Equal(t, int64(1), resp.CourierID)
	assert.Equal(t, "unassigned", resp.Status)
}

func TestDeliveryUseCase_UnassignDelivery_NoOrderID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryUnassignRequest{
		OrderID: "",
	}

	resp, err := uc.UnassignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrNoOrderID, err)
	assert.Equal(t, model.DeliveryUnassignResponse{}, resp)
}

func TestDeliveryUseCase_UnassignDelivery_OrderNotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryUnassignRequest{
		OrderID: "NONEXISTENT",
	}

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockDeliveryRepo.EXPECT().
		CouriersDelivery(ctx, "NONEXISTENT").
		Return(nil, repository.ErrOrderIDNotFound)

	resp, err := uc.UnassignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderIDNotFound, err)
	assert.Equal(t, model.DeliveryUnassignResponse{}, resp)
}

func TestDeliveryUseCase_UnassignDelivery_DeleteError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryUnassignRequest{
		OrderID: "ORDER123",
	}

	deliveryDB := &model.DeliveryDB{
		ID:        1,
		CourierID: 1,
		OrderID:   "ORDER123",
	}

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockDeliveryRepo.EXPECT().
		CouriersDelivery(ctx, "ORDER123").
		Return(deliveryDB, nil)

	mockDeliveryRepo.EXPECT().
		DeleteDelivery(ctx, "ORDER123").
		Return(repository.ErrOrderIDNotFound)

	resp, err := uc.UnassignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderIDNotFound, err)
	assert.Equal(t, model.DeliveryUnassignResponse{}, resp)
}

func TestDeliveryUseCase_UnassignDelivery_CourierNotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	mockDeliveryRepo := mocks.NewMockdeliveryRepository(ctrl)
	mockTxRunner := mocks.NewMocktxRunner(ctrl)

	uc := NewDelieveryUseCase(mockCourierRepo, mockDeliveryRepo, mockTxRunner)

	ctx := context.Background()
	req := &model.DeliveryUnassignRequest{
		OrderID: "ORDER123",
	}

	deliveryDB := &model.DeliveryDB{
		ID:        1,
		CourierID: 999,
		OrderID:   "ORDER123",
	}

	expectedErr := repository.ErrCourierNotFound

	mockTxRunner.EXPECT().
		Run(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockDeliveryRepo.EXPECT().
		CouriersDelivery(ctx, "ORDER123").
		Return(deliveryDB, nil)

	mockDeliveryRepo.EXPECT().
		DeleteDelivery(ctx, "ORDER123").
		Return(nil)

	mockCourierRepo.EXPECT().
		GetCourierById(ctx, int64(999)).
		Return(nil, expectedErr)

	resp, err := uc.UnassignDelivery(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, model.DeliveryUnassignResponse{}, resp)
}

func TestTransportTypeTime(t *testing.T) {
	testCases := []struct {
		name          string
		transportType string
		expected      time.Duration
		expectError   bool
	}{
		{"car", "car", 20 * time.Second, false},
		{"scooter", "scooter", 40 * time.Second, false},
		{"on_foot", "on_foot", 60 * time.Second, false},
		{"invalid", "rocket", 0, true},
		{"empty", "", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestCheckFreeCouriers_Success(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockCourierRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockCourierRepo.EXPECT().
		FreeCouriersWithInterval(gomock.Any()).
		Return(nil).
		MinTimes(2)

	go CheckFreeCouriersWithInterval(ctx, uc, 50*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	cancel()

	time.Sleep(50 * time.Millisecond)
}

func TestCheckFreeCouriers_RepositoryError(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockCourierRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockCourierRepo.EXPECT().
		FreeCouriersWithInterval(gomock.Any()).
		Return(errors.New("database error")).
		MinTimes(2)

	go CheckFreeCouriersWithInterval(ctx, uc, 50*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	cancel()

	time.Sleep(50 * time.Millisecond)
}

func TestCheckFreeCouriers_ContextCancellation(t *testing.T) {
	defer goleak.VerifyNone(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourierRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockCourierRepo)

	ctx, cancel := context.WithCancel(context.Background())

	go CheckFreeCouriersWithInterval(ctx, uc, 50*time.Millisecond)

	cancel()

	time.Sleep(50 * time.Millisecond)

}
