package order

import (
	"context"
	"errors"
	"fmt"
	"courier-service/internal/model"
	completeUC "courier-service/internal/usecase/delivery/complete"
)

type orderChanged interface {
	HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error
}

type OrderChangedFactory struct {
	assign assignUseCase
	unassign unassignUseCase
	complete completeUseCase
}

func NewOrderChangedFactory(assign assignUseCase, unassign unassignUseCase, complete completeUseCase) *OrderChangedFactory {
	return &OrderChangedFactory{
		assign: assign,
		unassign: unassign,
		complete: complete,
	}
}

type createdOrderHandler struct {
	assign assignUseCase
}

type completedOrderHandler struct {
	complete completeUseCase
}

type cancelledOrderHandler struct {
	unassign unassignUseCase
}

func (c *createdOrderHandler) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	_, err := c.assign.Assign(ctx, orderIDtoDeliveryAssignRequest(orderID))
	if err != nil {
		return err
	}
	return nil
}

func (c *completedOrderHandler) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	err := c.complete.Complete(ctx, orderIDtoCompleteDeliveryRequest(orderID))
	if err != nil {
		if errors.Is(err, completeUC.ErrOrderNotFound) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (c *cancelledOrderHandler) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	_, err := c.unassign.Unassign(ctx, orderIDtoDeliveryUnassignRequest(orderID))
	if err != nil {
		return err
	}
	return nil
}

func (f *OrderChangedFactory) GetOrderChanged(status model.OrderStatus) orderChanged {
	switch status {
	case model.OrderStatusCreated:
		return &createdOrderHandler{
			assign: f.assign,
		}
	case model.OrderStatusCompleted:
		return &completedOrderHandler{
			complete: f.complete,
		}
	case model.OrderStatusCancelled:
		return &cancelledOrderHandler{
			unassign: f.unassign,
		}
	default:
		return nil
	}
}
