package handlers

import (
	"bytes"
	"context"
	"courier-service/internal/handlers/dto"
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

func TestCourierHandler_GetCourierById(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		routePattern   string
		url            string
		prepare        func(uc *mocks.MockсourierUseCase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name:         "success",
			routePattern: "/courier/{id}",
			url:          "/courier/123",
			prepare: func(uc *mocks.MockсourierUseCase) {
				expectedCourier := model.Courier{
					ID:            123,
					Name:          "John Doe",
					Phone:         "+79991234567",
					Status:        "available",
					TransportType: "car",
				}

				uc.EXPECT().
					GetCourierById(gomock.Any(), int64(123)).
					Return(expectedCourier, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result model.Courier
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, int64(123), result.ID)
				assert.Equal(t, "John Doe", result.Name)
			},
		},
		{
			name:         "not found",
			routePattern: "/courier/{id}",
			url:          "/courier/999",
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					GetCourierById(gomock.Any(), int64(999)).
					Return(model.Courier{}, usecase.ErrCourierNotFound)
			},
			wantStatusCode: http.StatusNotFound,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// при необходимости можно проверить тело ответа/сообщение об ошибке
			},
		},
		{
			name:         "invalid id",
			routePattern: "/courier/{id}",
			url:          "/courier/invalid",
			prepare: func(uc *mocks.MockсourierUseCase) {
				// usecase вызываться не должен
			},
			wantStatusCode: http.StatusBadRequest,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// можно проверить сообщение об ошибке, если нужно
			},
		},
		{
			name:         "internal error",
			routePattern: "/courier/{id}",
			url:          "/courier/1",
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					GetCourierById(gomock.Any(), int64(1)).
					Return(model.Courier{}, assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// можно проверить тело с ошибкой
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockсourierUseCase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUseCase)
			}

			controller := NewCourierController(mockUseCase)

			router := chi.NewRouter()
			router.Get(tt.routePattern, controller.GetCourierById)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.expectations != nil {
				tt.expectations(t, rr)
			}
		})
	}
}

func TestCourierHandler_GetAllCouriers(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		url            string
		prepare        func(uc *mocks.MockсourierUseCase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name: "success",
			url:  "/courier",
			prepare: func(uc *mocks.MockсourierUseCase) {
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

				uc.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return(expectedCouriers, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result []model.Courier
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Len(t, result, 2)
				assert.Equal(t, int64(1), result[0].ID)
				assert.Equal(t, int64(2), result[1].ID)
			},
		},
		{
			name: "internal error",
			url:  "/couriers",
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return(nil, assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// при желании можно проверить тело/сообщение об ошибке
			},
		},
		{
			name: "empty list",
			url:  "/couriers",
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					GetAllCouriers(gomock.Any()).
					Return([]model.Courier{}, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result []model.Courier
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockсourierUseCase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUseCase)
			}

			controller := NewCourierController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rr := httptest.NewRecorder()

			controller.GetAllCouriers(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.expectations != nil {
				tt.expectations(t, rr)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestCourierHandler_UpdateCourier(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		url            string
		requestBody    []byte
		prepare        func(uc *mocks.MockсourierUseCase)
		wantStatusCode int
		expectations   expectationsFn
	}{
		{
			name: "missing required fields",
			url:  "/courier/",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:   1,
					Name: stringPtr(""),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(usecase.ErrInvalidUpdate)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "unknown transport type",
			url:  "/courier/",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:            1,
					TransportType: stringPtr("spaceship"),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(usecase.ErrUnknownTransportType)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing id",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					Name: stringPtr("Updated Name"),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare:        nil, // usecase не вызывается
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			url:            "/courier",
			requestBody:    []byte("invalid json"),
			prepare:        nil, // usecase не вызывается
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "courier not found",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:   999,
					Name: stringPtr("Updated Name"),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(usecase.ErrCourierNotFound)
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "phone exists",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:    1,
					Phone: stringPtr("+79991234567"),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(usecase.ErrPhoneNumberExists)
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "invalid phone",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:    1,
					Phone: stringPtr("invalid"),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(usecase.ErrInvalidPhoneNumber)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:   1,
					Name: stringPtr("Updated Name"),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierUpdateRequestDTO{
					ID:            1,
					Name:          stringPtr("Updated Name"),
					TransportType: stringPtr(string(model.TransportTypeCar)),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					UpdateCourier(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, c model.Courier) error {
						assert.Equal(t, int64(1), c.ID)
						assert.Equal(t, "Updated Name", c.Name)
						assert.Equal(t, model.TransportTypeCar, c.TransportType)
						return nil
					})
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockсourierUseCase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUseCase)
			}

			controller := NewCourierController(mockUseCase)

			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewReader(tt.requestBody))
			rr := httptest.NewRecorder()

			controller.UpdateCourier(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.expectations != nil {
				tt.expectations(t, rr)
			}
		})
	}
}

func TestCourierHandler_CreateCourier(t *testing.T) {
	type expectationsFn func(t *testing.T, rr *httptest.ResponseRecorder)

	tests := []struct {
		name           string
		url            string
		requestBody    []byte
		prepare        func(uc *mocks.MockсourierUseCase)
		expectations   expectationsFn
		wantStatusCode int
	}{
		{
			name: "success",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierCreateRequestDTO{
					Name:          "John Doe",
					Phone:         "+79991234567",
					TransportType: string(model.TransportTypeCar),
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
			},
			wantStatusCode: http.StatusCreated,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

				var result map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, "1", result["id"])
				assert.Equal(t, "Courier created successfully", result["message"])
			},
		},
		{
			name:           "invalid json",
			url:            "/courier",
			requestBody:    []byte("invalid json"),
			prepare:        nil, // usecase вызываться не должен
			wantStatusCode: http.StatusBadRequest,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// при желании можно проверить тело/сообщение об ошибке
			},
		},
		{
			name: "phone exists",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierCreateRequestDTO{
					Name:          "John Doe",
					Phone:         "+79991234567",
					TransportType: "car",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(0), usecase.ErrPhoneNumberExists)
			},
			wantStatusCode: http.StatusConflict,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				// можно добавить проверку тела с текстом ошибки, если он есть
			},
		},
		{
			name: "validation error (invalid phone)",
			url:  "/couriers",
			requestBody: func() []byte {
				reqBody := dto.CourierCreateRequestDTO{
					Name:          "John Doe",
					Phone:         "invalid_phone",
					TransportType: "car",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(0), usecase.ErrInvalidPhoneNumber)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			url:  "/couriers",
			requestBody: func() []byte {
				reqBody := dto.CourierCreateRequestDTO{
					Name:          "John Doe",
					Phone:         "+79991234567",
					TransportType: "car",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(0), assert.AnError)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "missing required fields",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierCreateRequestDTO{
					Name:          "",
					Phone:         "+79991234567",
					TransportType: "car",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(0), usecase.ErrInvalidCreate)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "unknown transport type",
			url:  "/courier",
			requestBody: func() []byte {
				reqBody := dto.CourierCreateRequestDTO{
					Name:          "John Doe",
					Phone:         "+79991234567",
					TransportType: "airplane",
				}
				b, _ := json.Marshal(reqBody)
				return b
			}(),
			prepare: func(uc *mocks.MockсourierUseCase) {
				uc.EXPECT().
					CreateCourier(gomock.Any(), gomock.Any()).
					Return(int64(0), usecase.ErrUnknownTransportType)
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := mocks.NewMockсourierUseCase(ctrl)
			if tt.prepare != nil {
				tt.prepare(mockUseCase)
			}

			controller := NewCourierController(mockUseCase)

			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(tt.requestBody))
			rr := httptest.NewRecorder()

			controller.CreateCourier(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)

			if tt.expectations != nil {
				tt.expectations(t, rr)
			}
		})
	}
}
