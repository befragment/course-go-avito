package handlers

import (
	"bytes"
	"context"
	"courier-service/internal/handlers/mocks"
	"courier-service/internal/model"
	"courier-service/internal/usecase"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourierHandler_CreateCourier_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		TransportType: "car",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(int64(1), nil)

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

	var result map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "1", result["id"])
	assert.Equal(t, "Courier created successfully", result["message"])
}

func TestCourierHandler_CreateCourier_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	request := httptest.NewRequest("POST", "/courier", bytes.NewReader([]byte("invalid json")))
	response := httptest.NewRecorder()

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_CreateCourier_PhoneExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		TransportType: "car",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(int64(0), usecase.ErrPhoneNumberExists)

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestCourierHandler_GetCourierById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	router := chi.NewRouter()
	router.Get("/courier/{id}", controller.GetCourierById)

	request := httptest.NewRequest("GET", "/courier/123", nil)
	response := httptest.NewRecorder()

	expectedCourier := &model.Courier{
		ID:            123,
		Name:          "John Doe",
		Phone:         "+79991234567",
		Status:        "available",
		TransportType: "car",
	}
	mockUseCase.EXPECT().
		GetCourierById(gomock.Any(), int64(123)).
		Return(expectedCourier, nil)

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var result model.Courier
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, int64(123), result.ID)
	assert.Equal(t, "John Doe", result.Name)
}

func TestCourierHandler_GetCourierById_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	router := chi.NewRouter()
	router.Get("/courier/{id}", controller.GetCourierById)

	request := httptest.NewRequest("GET", "/courier/999", nil)
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		GetCourierById(gomock.Any(), int64(999)).
		Return(nil, usecase.ErrCourierNotFound)

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestCourierHandler_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	request := httptest.NewRequest("GET", "/courier", nil)
	response := httptest.NewRecorder()

	expectedCouriers := []model.Courier{
		{
			ID:            1,
			Name:          "John",
			Phone:         "+79991111111",
			Status:        "available",
			TransportType: "car",
		},
		{
			ID:            2,
			Name:          "Jane",
			Phone:         "+79992222222",
			Status:        "busy",
			TransportType: "bike",
		},
	}
	mockUseCase.EXPECT().
		GetAllCouriers(gomock.Any()).
		Return(expectedCouriers, nil)

	controller.GetAllCouriers(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var result []model.Courier
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, int64(2), result[1].ID)
}

func TestCourierHandler_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:            1,
		Name:          stringPtr("Updated Name"),
		TransportType: stringPtr("car"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, req *model.CourierUpdateRequest) error {
			assert.Equal(t, int64(1), req.ID)
			assert.Equal(t, "Updated Name", *req.Name)
			return nil
		})

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestCourierHandler_GetCourierById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	router := chi.NewRouter()
	router.Get("/couriers/{id}", controller.GetCourierById)

	request := httptest.NewRequest("GET", "/couriers/invalid", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_GetCourierById_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	router := chi.NewRouter()
	router.Get("/couriers/{id}", controller.GetCourierById)

	request := httptest.NewRequest("GET", "/couriers/1", nil)
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		GetCourierById(gomock.Any(), int64(1)).
		Return(nil, assert.AnError)

	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestCourierHandler_GetAllCouriers_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	request := httptest.NewRequest("GET", "/couriers", nil)
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		GetAllCouriers(gomock.Any()).
		Return(nil, assert.AnError)

	controller.GetAllCouriers(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestCourierHandler_GetAllCouriers_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	request := httptest.NewRequest("GET", "/couriers", nil)
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		GetAllCouriers(gomock.Any()).
		Return([]model.Courier{}, nil)

	controller.GetAllCouriers(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var result []model.Courier
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestCourierHandler_CreateCourier_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "invalid_phone",
		TransportType: "car",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/couriers", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(int64(0), usecase.ErrInvalidPhoneNumber)

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_CreateCourier_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		TransportType: "car",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/couriers", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(int64(0), assert.AnError)

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestCourierHandler_UpdateCourier_MissingID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		Name: stringPtr("Updated Name"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_UpdateCourier_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader([]byte("invalid json")))
	response := httptest.NewRecorder()

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_UpdateCourier_CourierNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:   999,
		Name: stringPtr("Updated Name"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(usecase.ErrCourierNotFound)

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestCourierHandler_UpdateCourier_PhoneExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:    1,
		Phone: stringPtr("+79991234567"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(usecase.ErrPhoneNumberExists)

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestCourierHandler_UpdateCourier_InvalidPhone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:    1,
		Phone: stringPtr("invalid"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(usecase.ErrInvalidPhoneNumber)

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_UpdateCourier_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:   1,
		Name: stringPtr("Updated Name"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()
	
	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(assert.AnError)

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestCourierHandler_CreateCourier_MissingRequiredFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierCreateRequest{
		Name:          "",
		Phone:         "+79991234567",
		TransportType: "car",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(int64(0), usecase.ErrInvalidCreate)

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_CreateCourier_UnknownTransportType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierCreateRequest{
		Name:          "John Doe",
		Phone:         "+79991234567",
		TransportType: "airplane",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/courier", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateCourier(gomock.Any(), gomock.Any()).
		Return(int64(0), usecase.ErrUnknownTransportType)

	controller.CreateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_UpdateCourier_MissingRequiredFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:   1,
		Name: stringPtr(""),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier/", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(usecase.ErrInvalidUpdate)

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCourierHandler_UpdateCourier_UnknownTransportType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockсourierUseCase(ctrl)
	controller := NewCourierController(mockUseCase)

	reqBody := model.CourierUpdateRequest{
		ID:            1,
		TransportType: stringPtr("spaceship"),
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("PUT", "/courier/", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UpdateCourier(gomock.Any(), gomock.Any()).
		Return(usecase.ErrUnknownTransportType)

	controller.UpdateCourier(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func stringPtr(s string) *string {
	return &s
}
