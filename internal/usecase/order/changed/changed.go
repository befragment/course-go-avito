package changed

import (
	"context"

	"courier-service/internal/model"
	logger "courier-service/pkg/logger"
)

type OrderChangedUseCase struct {
	factory      orderChangedFactory
	orderGateway orderGateway
	logger       logger.LoggerInterface
}

func NewOrderChangedUseCase(factory orderChangedFactory, gateway orderGateway, log logger.LoggerInterface) *OrderChangedUseCase {
	return &OrderChangedUseCase{
		factory:      factory,
		orderGateway: gateway,
		logger:       log,
	}
}

func (uc *OrderChangedUseCase) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	if status != model.OrderStatusCompleted {
		uc.logger.Debugf("sending grpc request for checking status for order %s", orderID)
		order, err := uc.orderGateway.GetOrderById(ctx, orderID)
		if err != nil {
			return err
		}

		if order.Status != status {
			uc.logger.Warnf("order status mismatch: expected %s, got %s for order %s", status, order.Status, orderID)
			return ErrOrderStatusMismatch
		}
	}


	processor, ok := uc.factory.Get(status)
	if !ok {
		return nil
	}
	return processor.HandleOrderStatusChanged(ctx, status, orderID)
}
