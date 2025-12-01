package usecase

import (
	"context"
	"errors"
	"testing"

	"courier-service/internal/model"
	"courier-service/internal/repository"
	"courier-service/internal/usecase/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCourierUseCase_GetById_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	courierDB := &model.CourierDB{
		ID:            1,
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	mockRepo.EXPECT().
		GetCourierById(ctx, int64(1)).
		Return(courierDB, nil)

	result, err := uc.GetCourierById(ctx, int64(1))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "+79991234567", result.Phone)
}

func TestCourierUseCase_GetById_NotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().
		GetCourierById(ctx, int64(999)).
		Return(nil, repository.ErrCourierNotFound)

	result, err := uc.GetCourierById(ctx, int64(999))

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrCourierNotFound, err)
}

func TestCourierUseCase_GetAll_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	couriersDB := []model.CourierDB{
		{ID: int64(1), Name: "John", Phone: "+79991234567", Status: "available", TransportType: "car"},
		{ID: int64(2), Name: "Jane", Phone: "+79991234568", Status: "busy", TransportType: "scooter"},
	}

	mockRepo.EXPECT().
		GetAllCouriers(ctx).
		Return(couriersDB, nil)

	result, err := uc.GetAllCouriers(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, int64(2), result[1].ID)
}

func TestCourierUseCase_GetAll_EmptyList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().
		GetAllCouriers(ctx).
		Return([]model.CourierDB{}, nil)

	result, err := uc.GetAllCouriers(ctx)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestCourierUseCase_Create_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	req := &model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	mockRepo.EXPECT().
		ExistsCourierByPhone(ctx, req.Phone).
		Return(false, nil)

	mockRepo.EXPECT().
		CreateCourier(ctx, gomock.Any()).
		Return(int64(1), nil)

	id, err := uc.CreateCourier(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestCourierUseCase_Create_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()

	testCases := []struct {
		name string
		req  *model.CourierCreateRequest
	}{
		{
			name: "missing name",
			req: &model.CourierCreateRequest{
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
		},
		{
			name: "missing phone",
			req: &model.CourierCreateRequest{
				Name:          "John",
				Status:        "available",
				TransportType: "car",
			},
		},
		{
			name: "missing status",
			req: &model.CourierCreateRequest{
				Name:          "John",
				Phone:         "+79991234567",
				TransportType: "car",
			},
		},
		{
			name: "missing transport_type",
			req: &model.CourierCreateRequest{
				Name:   "John",
				Phone:  "+79991234567",
				Status: "available",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := uc.CreateCourier(ctx, tc.req)

			assert.Error(t, err)
			assert.Equal(t, int64(0), id)
			assert.Equal(t, ErrInvalidCreate, err)
		})
	}
}

func TestCourierUseCase_Create_InvalidTransportType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	req := &model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "airplane",
	}

	id, err := uc.CreateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, ErrUnknownTransportType, err)
}

func TestCourierUseCase_Create_InvalidPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	req := &model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "123",
		Status:        "available",
		TransportType: "car",
	}

	id, err := uc.CreateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, ErrInvalidPhoneNumber, err)
}

func TestCourierUseCase_Create_PhoneExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	req := &model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}

	mockRepo.EXPECT().
		ExistsCourierByPhone(ctx, req.Phone).
		Return(true, nil)

	id, err := uc.CreateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Equal(t, ErrPhoneNumberExists, err)
}

func TestCourierUseCase_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	name := "Jane Doe"
	req := &model.CourierUpdateRequest{
		ID:   1,
		Name: &name,
	}

	mockRepo.EXPECT().
		UpdateCourier(ctx, gomock.Any()).
		Return(nil)

	err := uc.UpdateCourier(ctx, req)

	assert.NoError(t, err)
}

func TestCourierUseCase_Update_NoFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	req := &model.CourierUpdateRequest{
		ID: 1,
	}

	err := uc.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUpdate, err)
}

func TestCourierUseCase_Update_InvalidTransportType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	transportType := "rocket"
	req := &model.CourierUpdateRequest{
		ID:            1,
		TransportType: &transportType,
	}

	err := uc.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrUnknownTransportType, err)
}

func TestCourierUseCase_Update_InvalidPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	phone := "invalid"
	req := &model.CourierUpdateRequest{
		ID:    1,
		Phone: &phone,
	}

	err := uc.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPhoneNumber, err)
}

func TestCourierUseCase_Update_PhoneExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	phone := "+79991234567"
	req := &model.CourierUpdateRequest{
		ID:    1,
		Phone: &phone,
	}

	mockRepo.EXPECT().
		ExistsCourierByPhone(ctx, phone).
		Return(true, nil)

	err := uc.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrPhoneNumberExists, err)
}

func TestCourierUseCase_Update_CourierNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	name := "Jane Doe"
	req := &model.CourierUpdateRequest{
		ID:   999,
		Name: &name,
	}

	mockRepo.EXPECT().
		UpdateCourier(ctx, gomock.Any()).
		Return(repository.ErrCourierNotFound)

	err := uc.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrCourierNotFound, err)
}

func TestCourierUseCase_Update_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockсourierRepository(ctrl)
	uc := NewCourierUseCase(mockRepo)

	ctx := context.Background()
	name := "Jane Doe"
	req := &model.CourierUpdateRequest{
		ID:   1,
		Name: &name,
	}

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().
		UpdateCourier(ctx, gomock.Any()).
		Return(expectedErr)

	err := uc.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestValidPhoneNumber(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidPhoneNumber(tc.phone)
			assert.Equal(t, tc.expected, result)
		})
	}
}
