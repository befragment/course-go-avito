//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package delivery

import (
	"context"

	assign "courier-service/internal/usecase/delivery/assign"
	unassign "courier-service/internal/usecase/delivery/unassign"
)

type assignUsecase interface {
	Assign(context.Context, assign.DeliveryAssignRequest) (assign.DeliveryAssignResponse, error)
}

type unassignUsecase interface {
	Unassign(context.Context, unassign.DeliveryUnassignRequest) (unassign.DeliveryUnassignResponse, error)
}
