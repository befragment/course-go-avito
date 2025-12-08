package handlers

import (
	"bytes"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"courier-service/internal/handlers/mocks"
	"courier-service/internal/handlers/dto"
	"courier-service/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryHandler_AssignDelivery(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		requestBody    []byte
		prepare        func(uc *mocks.MockdeliveryUseCase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name: "success",
			requestBody: func() []byte {
				reqBody := dto.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				expectedResponse := usecase.DeliveryAssignResponse{
					CourierID:     1,
					OrderID:       "550e8400-e29b-41d4-a716-446655440000",
					TransportType: "car",
				}

				uc.EXPECT().
					AssignDelivery(gomock.Any(), gomock.Any()).
					Return(expectedResponse, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result usecase.DeliveryAssignResponse
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, int64(1), result.CourierID)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.OrderID)
			},
		},
		{
			name:        "invalid json",
			requestBody: []byte("invalid json"),
			prepare:     nil, // usecase не вызывается
			wantStatusCode: http.StatusBadRequest,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// можно проверить тело/сообщение об ошибке при необходимости
			},
		},
		{
			name: "all couriers busy",
			requestBody: func() []byte {
				reqBody := dto.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					AssignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryAssignResponse{}, usecase.ErrCouriersBusy)
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "missing order id",
			requestBody: func() []byte {
				reqBody := dto.DeliveryAssignRequestDTO{
					OrderID: "",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					AssignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryAssignResponse{}, usecase.ErrNoOrderID)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "order id exists",
			requestBody: func() []byte {
				reqBody := dto.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					AssignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryAssignResponse{}, usecase.ErrOrderIDExists)
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "internal error",
			requestBody: func() []byte {
				reqBody := dto.DeliveryAssignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					AssignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryAssignResponse{}, assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUseCase)
			}

			controller := NewDeliveryController(mockUseCase)

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
		prepare        func(uc *mocks.MockdeliveryUseCase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name: "success",
			requestBody: func() []byte {
				reqBody := dto.DeliveryUnassignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				expectedResponse := usecase.DeliveryUnassignResponse{
					OrderID:   "550e8400-e29b-41d4-a716-446655440000",
					Status:    "unassigned",
					CourierID: 1,
				}

				uc.EXPECT().
					UnassignDelivery(gomock.Any(), gomock.Any()).
					Return(expectedResponse, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result usecase.DeliveryUnassignResponse
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.OrderID)
				assert.Equal(t, "unassigned", result.Status)
				assert.Equal(t, int64(1), result.CourierID)
			},
		},
		{
			name:        "invalid json",
			requestBody: []byte("invalid json"),
			prepare:     nil, // usecase не вызывается
			wantStatusCode: http.StatusBadRequest,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// при необходимости можно проверить тело/сообщение об ошибке
			},
		},
		{
			name: "missing order id",
			requestBody: func() []byte {
				reqBody := dto.DeliveryUnassignRequestDTO{
					OrderID: "",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					UnassignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryUnassignResponse{}, usecase.ErrNoOrderID)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "order id not found",
			requestBody: func() []byte {
				reqBody := dto.DeliveryUnassignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					UnassignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryUnassignResponse{}, usecase.ErrOrderIDNotFound)
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			requestBody: func() []byte {
				reqBody := dto.DeliveryUnassignRequestDTO{
					OrderID: "550e8400-e29b-41d4-a716-446655440000",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockdeliveryUseCase) {
				uc.EXPECT().
					UnassignDelivery(gomock.Any(), gomock.Any()).
					Return(usecase.DeliveryUnassignResponse{}, assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockdeliveryUseCase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUseCase)
			}

			controller := NewDeliveryController(mockUseCase)

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