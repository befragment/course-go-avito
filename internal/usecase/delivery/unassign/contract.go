//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package unassign

import (
	"context"

	"courier-service/internal/model"
)

type courierRepository interface {
	GetCourierById(ctx context.Context, id int64) (model.Courier, error)
	UpdateCourier(ctx context.Context, courier model.Courier) error
}

type deliveryRepository interface {
	CreateDelivery(ctx context.Context, delivery model.Delivery) (model.Delivery, error)
	CouriersDelivery(ctx context.Context, orderID string) (model.Delivery, error)
	DeleteDelivery(ctx context.Context, orderID string) error
}

type txRunner interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
