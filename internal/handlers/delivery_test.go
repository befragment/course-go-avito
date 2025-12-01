package handlers

import (
	"bytes"
	"courier-service/internal/handlers/mocks"
	"courier-service/internal/model"
	"courier-service/internal/usecase"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryHandler_AssignDelivery_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryAssignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	expectedResponse := model.DeliveryAssignResponse{
		CourierID:     1,
		OrderID:       "550e8400-e29b-41d4-a716-446655440000",
		TransportType: "car",
		Deadline:      time.Now().Add(2 * time.Hour),
	}

	mockUseCase.EXPECT().
		AssignDelivery(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	controller.AssignDelivery(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var result model.DeliveryAssignResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.CourierID)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.OrderID)
}

func TestDeliveryHandler_AssignDelivery_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	request := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader([]byte("invalid json")))
	response := httptest.NewRecorder()

	controller.AssignDelivery(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestDeliveryHandler_AssignDelivery_AllCouriersBusy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryAssignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		AssignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryAssignResponse{}, usecase.ErrCouriersBusy)

	controller.AssignDelivery(response, request)

	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestDeliveryHandler_AssignDelivery_MissingOrderID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryAssignRequest{
		OrderID: "",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		AssignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryAssignResponse{}, usecase.ErrNoOrderID)

	controller.AssignDelivery(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestDeliveryHandler_AssignDelivery_OrderIDExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryAssignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		AssignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryAssignResponse{}, usecase.ErrOrderIDExists)

	controller.AssignDelivery(response, request)

	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestDeliveryHandler_AssignDelivery_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryAssignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/assign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		AssignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryAssignResponse{}, assert.AnError)

	controller.AssignDelivery(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestDeliveryHandler_UnassignDelivery_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryUnassignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	expectedResponse := model.DeliveryUnassignResponse{
		OrderID:   "550e8400-e29b-41d4-a716-446655440000",
		Status:    "unassigned",
		CourierID: 1,
	}

	mockUseCase.EXPECT().
		UnassignDelivery(gomock.Any(), gomock.Any()).
		Return(expectedResponse, nil)

	controller.UnassignDelivery(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	var result model.DeliveryUnassignResponse
	err := json.Unmarshal(response.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.OrderID)
	assert.Equal(t, "unassigned", result.Status)
}

func TestDeliveryHandler_UnassignDelivery_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	request := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader([]byte("invalid json")))
	response := httptest.NewRecorder()

	controller.UnassignDelivery(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestDeliveryHandler_UnassignDelivery_MissingOrderID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryUnassignRequest{
		OrderID: "",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UnassignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryUnassignResponse{}, usecase.ErrNoOrderID)

	controller.UnassignDelivery(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestDeliveryHandler_UnassignDelivery_OrderIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryUnassignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UnassignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryUnassignResponse{}, usecase.ErrOrderIDNotFound)

	controller.UnassignDelivery(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestDeliveryHandler_UnassignDelivery_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
	controller := NewDeliveryController(mockUseCase)

	reqBody := model.DeliveryUnassignRequest{
		OrderID: "550e8400-e29b-41d4-a716-446655440000",
	}
	body, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewReader(body))
	response := httptest.NewRecorder()

	mockUseCase.EXPECT().
		UnassignDelivery(gomock.Any(), gomock.Any()).
		Return(model.DeliveryUnassignResponse{}, assert.AnError)

	controller.UnassignDelivery(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
