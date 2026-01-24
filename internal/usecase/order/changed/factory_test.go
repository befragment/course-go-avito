package changed_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"courier-service/internal/model"
	"courier-service/internal/usecase/order/changed"
)

func TestFactory_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       model.OrderStatus
		processors   map[model.OrderStatus]changed.Processor
		expectations func(t *testing.T, processor changed.Processor, found bool)
	}{
		{
			name:   "success: processor found for created status",
			status: model.OrderStatusCreated,
			processors: func() map[model.OrderStatus]changed.Processor {
				ctrl := gomock.NewController(t)
				mockProcessor := NewMockProcessor(ctrl)
				return map[model.OrderStatus]changed.Processor{
					model.OrderStatusCreated: mockProcessor,
				}
			}(),
			expectations: func(t *testing.T, processor changed.Processor, found bool) {
				assert.True(t, found)
				assert.NotNil(t, processor)
			},
		},
		{
			name:   "success: processor found for cancelled status",
			status: model.OrderStatusCancelled,
			processors: func() map[model.OrderStatus]changed.Processor {
				ctrl := gomock.NewController(t)
				mockProcessor := NewMockProcessor(ctrl)
				return map[model.OrderStatus]changed.Processor{
					model.OrderStatusCancelled: mockProcessor,
				}
			}(),
			expectations: func(t *testing.T, processor changed.Processor, found bool) {
				assert.True(t, found)
				assert.NotNil(t, processor)
			},
		},
		{
			name:   "success: processor found for completed status",
			status: model.OrderStatusCompleted,
			processors: func() map[model.OrderStatus]changed.Processor {
				ctrl := gomock.NewController(t)
				mockProcessor := NewMockProcessor(ctrl)
				return map[model.OrderStatus]changed.Processor{
					model.OrderStatusCompleted: mockProcessor,
				}
			}(),
			expectations: func(t *testing.T, processor changed.Processor, found bool) {
				assert.True(t, found)
				assert.NotNil(t, processor)
			},
		},
		{
			name:   "success: processor not found for status",
			status: model.OrderStatus("in_progress"),
			processors: func() map[model.OrderStatus]changed.Processor {
				ctrl := gomock.NewController(t)
				mockProcessor := NewMockProcessor(ctrl)
				return map[model.OrderStatus]changed.Processor{
					model.OrderStatusCreated: mockProcessor,
				}
			}(),
			expectations: func(t *testing.T, processor changed.Processor, found bool) {
				assert.False(t, found)
				assert.Nil(t, processor)
			},
		},
		{
			name:       "success: empty factory",
			status:     model.OrderStatusCreated,
			processors: map[model.OrderStatus]changed.Processor{},
			expectations: func(t *testing.T, processor changed.Processor, found bool) {
				assert.False(t, found)
				assert.Nil(t, processor)
			},
		},
		{
			name:   "success: multiple processors registered",
			status: model.OrderStatusCancelled,
			processors: func() map[model.OrderStatus]changed.Processor {
				ctrl := gomock.NewController(t)
				mockProcessor1 := NewMockProcessor(ctrl)
				mockProcessor2 := NewMockProcessor(ctrl)
				mockProcessor3 := NewMockProcessor(ctrl)
				return map[model.OrderStatus]changed.Processor{
					model.OrderStatusCreated:   mockProcessor1,
					model.OrderStatusCancelled: mockProcessor2,
					model.OrderStatusCompleted: mockProcessor3,
				}
			}(),
			expectations: func(t *testing.T, processor changed.Processor, found bool) {
				assert.True(t, found)
				assert.NotNil(t, processor)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			factory := changed.NewFactory(tc.processors)

			processor, found := factory.Get(tc.status)

			if tc.expectations != nil {
				tc.expectations(t, processor, found)
			}
		})
	}
}
