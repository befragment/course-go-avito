package courier_test

import (
	"reflect"
	"time"

	"github.com/golang/mock/gomock"
)

// MockDeliveryCalculator is a mock of DeliveryCalculator interface.
type MockDeliveryCalculator struct {
	ctrl     *gomock.Controller
	recorder *MockDeliveryCalculatorMockRecorder
}

// MockDeliveryCalculatorMockRecorder is the mock recorder for MockDeliveryCalculator.
type MockDeliveryCalculatorMockRecorder struct {
	mock *MockDeliveryCalculator
}

// NewMockDeliveryCalculator creates a new mock instance.
func NewMockDeliveryCalculator(ctrl *gomock.Controller) *MockDeliveryCalculator {
	mock := &MockDeliveryCalculator{ctrl: ctrl}
	mock.recorder = &MockDeliveryCalculatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeliveryCalculator) EXPECT() *MockDeliveryCalculatorMockRecorder {
	return m.recorder
}

// CalculateDeadline mocks base method.
func (m *MockDeliveryCalculator) CalculateDeadline() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateDeadline")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// CalculateDeadline indicates an expected call of CalculateDeadline.
func (mr *MockDeliveryCalculatorMockRecorder) CalculateDeadline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateDeadline", reflect.TypeOf((*MockDeliveryCalculator)(nil).CalculateDeadline))
}
