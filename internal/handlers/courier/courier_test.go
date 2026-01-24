package courier_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"courier-service/internal/handlers/courier"
	"courier-service/internal/model"
	usecase "courier-service/internal/usecase/courier"
)

func TestCourierHandler_GetById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		courierID      string
		prepare        func(courierUC *MockcourierUseCase)
		wantStatusCode int
		expectations   func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name:      "success: found courier with valid id",
			courierID: "1",
			prepare: func(courierUC *MockcourierUseCase) {
				courierUC.EXPECT().
					GetCourierById(gomock.Any(), int64(1)).
					Return(model.Courier{
						ID:            1,
						Name:          "John Doe",
						Phone:         "+79991234567",
						Status:        "available",
						TransportType: "car",
					}, nil)
			},
			wantStatusCode: http.StatusOK,
			expectations: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var result model.Courier
				err := json.Unmarshal(rr.Body.Bytes(), &result)
				require.NoError(t, err)

				assert.Equal(t, int64(1), result.ID)
				assert.Equal(t, "John Doe", result.Name)
				assert.Equal(t, "+79991234567", result.Phone)
				assert.Equal(t, model.CourierStatus("available"), result.Status)
				assert.Equal(t, model.CourierTransportType("car"), result.TransportType)
			},
		},
		{
			name:      "error: courier not found",
			courierID: "999",
			prepare: func(courierUC *MockcourierUseCase) {
				courierUC.EXPECT().
					GetCourierById(gomock.Any(), int64(999)).
					Return(model.Courier{}, usecase.ErrCourierNotFound)
			},
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUseCase := NewMockcourierUseCase(ctrl)
			if tc.prepare != nil {
				tc.prepare(mockUseCase)
			}

			controller := courier.NewCourierController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/courier/"+tc.courierID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.courierID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			controller.GetCourierById(rr, req)

			assert.Equal(t, tc.wantStatusCode, rr.Code)

			// 9. Дополнительные проверки
			if tc.expectations != nil {
				tc.expectations(t, rr)
			}
		})
	}
}
