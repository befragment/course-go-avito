package changed

import (
    "context"
    "courier-service/internal/model"
)

type Processor interface {
    HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error
}

type orderChangedFactory interface {
    Get(status model.OrderStatus) (Processor, bool)
}

type orderGateway interface {
    GetOrderById(ctx context.Context, orderID string) (model.Order, error)
}