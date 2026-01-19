//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package courier

import (
	"context"
	"courier-service/internal/model"
)

type courierUseCase interface {
	GetCourierById(ctx context.Context, id int64) (model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, courier model.Courier) (int64, error)
	UpdateCourier(ctx context.Context, courier model.Courier) error
}
