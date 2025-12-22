//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package delivery

import (
	"context"

	assign "courier-service/internal/usecase/delivery/assign"
)

type assignUsecase interface {
	Assign(context.Context, string) (assign.DeliveryAssignResponse, error)
}

type unassignUsecase interface {
	Unassign(context.Context, string) (int64, error)
}
