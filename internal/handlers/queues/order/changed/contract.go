package order

import (
	"context"
	"courier-service/internal/model"
)

type orderChangedUseCase interface {
	HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error
}
