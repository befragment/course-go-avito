//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package processor

import (
	"context"

	assign "courier-service/internal/usecase/delivery/assign"
)

type assignUseCase interface {
	Assign(ctx context.Context, OrderID string) (assign.DeliveryAssignResponse, error)
}

type unassignUseCase interface {
	Unassign(ctx context.Context, OrderID string) (int64, error)
}

type completeUseCase interface {
	Complete(ctx context.Context, OrderID string) error
}
