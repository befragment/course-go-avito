package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	deliveryhandler "courier-service/internal/handlers/delivery"
	assignusecase "courier-service/internal/usecase/delivery/assign"
	unassignusecase "courier-service/internal/usecase/delivery/unassign"
)

func TestDeliveryHandler_AssignDelivery(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		requestBody    []byte
		prepare        func(uc *MockassignUsecase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name: "success",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockassignUsecase) {
				expectedResponse := assignusecase.DeliveryAssignResponse{
					CourierID:     1,
					OrderID:       "550e8400-e29b-41d4-a716-446655440000",
					TransportType: "car",
				}

				uc.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(expectedResponse, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result deliveryhandler.DeliveryAssignResponseDTO
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, int64(1), result.CourierID)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.OrderID)
				assert.Equal(t, "car", result.TransportType)
			},
		},
		{
			name:           "invalid json",
			requestBody:    []byte("invalid json"),
			prepare:        nil, // usecase не вызывается
			wantStatusCode: http.StatusBadRequest,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
			},
		},
		{
			name: "all couriers busy",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockassignUsecase) {
				uc.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(assignusecase.DeliveryAssignResponse{}, assignusecase.ErrCouriersBusy)
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "missing order id",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryAssignRequestDTO{
					OrderID: "",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockassignUsecase) {
				uc.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(assignusecase.DeliveryAssignResponse{}, assignusecase.ErrNoOrderID)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "order id exists",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockassignUsecase) {
				uc.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(assignusecase.DeliveryAssignResponse{}, assignusecase.ErrOrderIDExists)
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "internal error",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockassignUsecase) {
				uc.EXPECT().
					Assign(gomock.Any(), gomock.Any()).
					Return(assignusecase.DeliveryAssignResponse{}, assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAssignUsecase := NewMockassignUsecase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockAssignUsecase)
			}

			controller := deliveryhandler.NewDeliveryController(mockAssignUsecase, nil)

			req := httptest.NewRequest(http.MethodPost, "/delivery/assign", bytes.NewReader(tt.requestBody))
			rr := httptest.NewRecorder()

			controller.AssignDelivery(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.expectations != nil {
				tt.expectations(t, rr)
			}
		})
	}
}

func TestDeliveryHandler_UnassignDelivery(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		requestBody    []byte
		prepare        func(uc *MockunassignUsecase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name: "success",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryUnassignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockunassignUsecase) {
				expectedResponse := int64(1)

				uc.EXPECT().
					Unassign(gomock.Any(), gomock.Any()).
					Return(expectedResponse, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result deliveryhandler.DeliveryUnassignResponseDTO
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, int64(1), result.CourierID)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.OrderID)
				assert.Equal(t, deliveryhandler.UnassignedStatus, result.Status)
			},
		},
		{
			name:           "invalid json",
			requestBody:    []byte("invalid json"),
			prepare:        nil, // usecase не вызывается
			wantStatusCode: http.StatusBadRequest,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// при необходимости можно проверить тело/сообщение об ошибке
			},
		},
		{
			name: "missing order id",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryUnassignRequestDTO{
					OrderID: "",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockunassignUsecase) {
				uc.EXPECT().
					Unassign(gomock.Any(), gomock.Any()).
					Return(int64(0), unassignusecase.ErrNoOrderID)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "order id not found",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryUnassignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockunassignUsecase) {
				uc.EXPECT().
					Unassign(gomock.Any(), gomock.Any()).
					Return(int64(0), unassignusecase.ErrOrderIDNotFound)
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			requestBody: func() []byte {
				reqBody := deliveryhandler.DeliveryUnassignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *MockunassignUsecase) {
				uc.EXPECT().
					Unassign(gomock.Any(), gomock.Any()).
					Return(int64(0), assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUnassignUsecase := NewMockunassignUsecase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUnassignUsecase)
			}

			controller := deliveryhandler.NewDeliveryController(nil, mockUnassignUsecase)

			req := httptest.NewRequest(http.MethodPost, "/delivery/unassign", bytes.NewReader(tt.requestBody))
			rr := httptest.NewRecorder()

			controller.UnassignDelivery(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.expectations != nil {
				tt.expectations(t, rr)
			}
		})
	}
}
