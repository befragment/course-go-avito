package changed

import (
	"context"

	"courier-service/internal/model"
	logger "courier-service/pkg/logger"
)

type OrderChangedUseCase struct {
	factory      orderChangedFactory
	orderGateway orderGateway
	logger       logger.Interface
}

func NewOrderChangedUseCase(factory orderChangedFactory, gateway orderGateway, log logger.Interface) *OrderChangedUseCase {
	return &OrderChangedUseCase{
		factory:      factory,
		orderGateway: gateway,
		logger:       log,
	}
}

func (uc *OrderChangedUseCase) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	order, err := uc.orderGateway.GetOrderById(ctx, orderID)
	if err != nil {
		return err
	}

	uc.logger.Debugf("sending grpc request for checking status for order %s", orderID)
	if order.Status != status {
		uc.logger.Warnf("order status mismatch: expected %s, got %s for order %s", status, order.Status, orderID)
		return ErrOrderStatusMismatch
	}

	processor, ok := uc.factory.Get(status)
	if !ok {
		return nil
	}
	return processor.HandleOrderStatusChanged(ctx, status, orderID)
}
