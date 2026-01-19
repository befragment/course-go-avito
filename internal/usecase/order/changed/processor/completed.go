package processor

import (
	"context"
	"errors"
	"fmt"

	"courier-service/internal/model"
	complete "courier-service/internal/usecase/delivery/complete"
	changed "courier-service/internal/usecase/order/changed"
)

type CompletedProcessor struct {
	completeUC completeUseCase
}

func NewCompletedProcessor(completeUC completeUseCase) *CompletedProcessor {
	return &CompletedProcessor{completeUC: completeUC}
}

func (p *CompletedProcessor) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	err := p.completeUC.Complete(ctx, orderID)
	if err != nil {
		if errors.Is(err, complete.ErrOrderNotFound) {
			return changed.ErrOrderNotFound
		}
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
