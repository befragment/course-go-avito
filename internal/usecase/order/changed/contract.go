//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package changed

import (
	"context"

	"courier-service/internal/model"
)

type Processor interface {
	HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error
}

type orderChangedFactory interface {
	Get(status model.OrderStatus) (Processor, bool)
}

type orderGateway interface {
	GetOrderById(ctx context.Context, orderID string) (model.Order, error)
}

type logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
