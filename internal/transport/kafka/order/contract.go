package order

import (
	"context"
	assign "courier-service/internal/usecase/delivery/assign"
	complete "courier-service/internal/usecase/delivery/complete"
	unassign "courier-service/internal/usecase/delivery/unassign"
)

type assignUseCase interface {
	Assign(ctx context.Context, req assign.DeliveryAssignRequest) (assign.DeliveryAssignResponse, error)
}

type unassignUseCase interface {
	Unassign(ctx context.Context, req unassign.DeliveryUnassignRequest) (unassign.DeliveryUnassignResponse, error)
}

type completeUseCase interface {
	Complete(ctx context.Context, req complete.CompleteDeliveryRequest) error
}