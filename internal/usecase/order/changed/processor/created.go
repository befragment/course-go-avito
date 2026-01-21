package processor

import (
	"context"

	"courier-service/internal/model"
)

type CreatedProcessor struct {
	assignUC assignUseCase
}

func NewCreatedProcessor(assignUC assignUseCase) *CreatedProcessor {
	return &CreatedProcessor{assignUC: assignUC}
}

func (p *CreatedProcessor) HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error {
	_, err := p.assignUC.Assign(ctx, orderID)
	return err
}
