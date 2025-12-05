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

func TestCourierUseCase_GetById(t *testing.T) {
	tests := []struct {
		name         string
		courierID    int64
		prepare      func(repo *mocks.MockсourierRepository)
		expectations func(t *testing.T, result *model.Courier, err error)
	}{
		{
			name:      "success: courier found",
			courierID: 1,
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					GetCourierById(gomock.Any(), int64(1)).
					Return(&repository.CourierDB{
						ID:            1,
						Name:          "John Doe",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)
			},
			expectations: func(t *testing.T, result *model.Courier, err error) {
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
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					GetCourierById(gomock.Any(), int64(999)).
					Return(nil, repository.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, result *model.Courier, err error) {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, ErrCourierNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockсourierRepository(ctrl)
			uc := NewCourierUseCase(mockRepo)

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
		prepare      func(repo *mocks.MockсourierRepository)
		expectations func(t *testing.T, result []model.Courier, err error)
	}{
		{
			name: "success: returns multiple couriers",
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return([]repository.CourierDB{
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
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return([]repository.CourierDB{}, nil)
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

			mockRepo := mocks.NewMockсourierRepository(ctrl)
			uc := NewCourierUseCase(mockRepo)

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
		request      *model.CourierCreateRequest
		prepare      func(repo *mocks.MockсourierRepository)
		expectations func(t *testing.T, id int64, err error)
	}{
		{
			name: "success: courier created",
			request: &model.CourierCreateRequest{
				Name:          "John Doe",
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
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
			request: &model.CourierCreateRequest{
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations - validation happens before any repo calls
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrInvalidCreate, err)
			},
		},
		{
			name: "error: missing phone",
			request: &model.CourierCreateRequest{
				Name:          "John",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrInvalidCreate, err)
			},
		},
		{
			name: "error: missing status",
			request: &model.CourierCreateRequest{
				Name:          "John",
				Phone:         "+79991234567",
				TransportType: "car",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrInvalidCreate, err)
			},
		},
		{
			name: "error: missing transport_type",
			request: &model.CourierCreateRequest{
				Name:   "John",
				Phone:  "+79991234567",
				Status: "available",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrInvalidCreate, err)
			},
		},
		{
			name: "error: invalid transport type",
			request: &model.CourierCreateRequest{
				Name:          "John Doe",
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "airplane",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrUnknownTransportType, err)
			},
		},
		{
			name: "error: invalid phone",
			request: &model.CourierCreateRequest{
				Name:          "John Doe",
				Phone:         "123",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrInvalidPhoneNumber, err)
			},
		},
		{
			name: "error: phone already exists",
			request: &model.CourierCreateRequest{
				Name:          "John Doe",
				Phone:         "+79991234567",
				Status:        "available",
				TransportType: "car",
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					ExistsCourierByPhone(gomock.Any(), "+79991234567").
					Return(true, nil)
			},
			expectations: func(t *testing.T, id int64, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), id)
				assert.Equal(t, ErrPhoneNumberExists, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockсourierRepository(ctrl)
			uc := NewCourierUseCase(mockRepo)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockRepo)
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
		request      *model.CourierUpdateRequest
		prepare      func(repo *mocks.MockсourierRepository)
		expectations func(t *testing.T, err error)
	}{
		{
			name: "success: courier updated",
			request: &model.CourierUpdateRequest{
				ID:   1,
				Name: &nameUpdate,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
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
			request: &model.CourierUpdateRequest{
				ID: 1,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations - validation happens before any repo calls
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidUpdate, err)
			},
		},
		{
			name: "error: invalid transport type",
			request: &model.CourierUpdateRequest{
				ID:            1,
				TransportType: &transportTypeUpdate,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrUnknownTransportType, err)
			},
		},
		{
			name: "error: invalid phone",
			request: &model.CourierUpdateRequest{
				ID:    1,
				Phone: &invalidPhone,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				// No mock expectations
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidPhoneNumber, err)
			},
		},
		{
			name: "error: phone already exists",
			request: &model.CourierUpdateRequest{
				ID:    1,
				Phone: &phoneUpdate,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					ExistsCourierByPhone(gomock.Any(), phoneUpdate).
					Return(true, nil)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrPhoneNumberExists, err)
			},
		},
		{
			name: "error: courier not found",
			request: &model.CourierUpdateRequest{
				ID:   999,
				Name: &nameUpdate,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
				repo.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(repository.ErrCourierNotFound)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrCourierNotFound, err)
			},
		},
		{
			name: "error: repository error",
			request: &model.CourierUpdateRequest{
				ID:   1,
				Name: &nameUpdate,
			},
			prepare: func(repo *mocks.MockсourierRepository) {
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

			mockRepo := mocks.NewMockсourierRepository(ctrl)
			uc := NewCourierUseCase(mockRepo)

			ctx := context.Background()

			if tc.prepare != nil {
				tc.prepare(mockRepo)
			}

			err := uc.UpdateCourier(ctx, tc.request)

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
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ValidPhoneNumber(tc.phone)
			assert.Equal(t, tc.expected, result)
		})
	}
}
