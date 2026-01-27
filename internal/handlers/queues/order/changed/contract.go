package order

import (
	"context"

	"courier-service/internal/model"
)

type orderChangedUseCase interface {
	HandleOrderStatusChanged(ctx context.Context, status model.OrderStatus, orderID string) error
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
