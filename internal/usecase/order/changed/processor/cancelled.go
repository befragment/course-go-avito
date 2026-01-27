package processor

import (
	"context"

	"courier-service/internal/model"
)

type CancelledProcessor struct {
	unassignUC unassignUseCase
}

func NewCancelledProcessor(unassignUC unassignUseCase) *CancelledProcessor {
	return &CancelledProcessor{unassignUC: unassignUC}
}

func (p *CancelledProcessor) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	_, err := p.unassignUC.Unassign(ctx, orderID)
	return err
}
