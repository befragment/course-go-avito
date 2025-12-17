package order

import (
	"context"
	"courier-service/internal/model"
	courierrepoerrors "courier-service/internal/repository/courier"
	deliveryrepoerrors "courier-service/internal/repository/delivery"
	"errors"
	"log"
	"time"
)

type OrderMonitoringUseCase struct {
	orderGateway              orderGateway
	courierRepository         courierRepository
	txRunner                  txRunner
	deliveryRepository        deliveryRepository
	deliveryCalculatorFactory deliveryCalculatorFactory
}

func NewOrderMonitoringUseCase(
	orderGateway orderGateway,
	courierRepository courierRepository,
	deliveryRepository deliveryRepository,
	txRunner txRunner,
	deliveryCalculatorFactory deliveryCalculatorFactory,
) *OrderMonitoringUseCase {
	return &OrderMonitoringUseCase{
		orderGateway:              orderGateway,
		courierRepository:         courierRepository,
		deliveryRepository:        deliveryRepository,
		txRunner:                  txRunner,
		deliveryCalculatorFactory: deliveryCalculatorFactory,
	}
}

func (u *OrderMonitoringUseCase) MonitorOrders(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			from := time.Now().Add(-interval)
			log.Printf("getting orders from gateway, cursor: %s\n", from.Format(time.RFC3339))
			orders, err := u.orderGateway.GetOrders(ctx, from)
			if err != nil {
				log.Printf("failed to get orders from gateway: %v\n", err)
				continue
			}
			for _, order := range orders {
				courier, _, err := createAssignmentInTransaction(
					ctx,
					u.txRunner,
					u.courierRepository,
					u.deliveryRepository,
					u.deliveryCalculatorFactory,
					order.ID,
				)
				if err != nil {
					log.Printf("failed to create assignment for order %s: %v\n", order.ID, err)
					continue
				}
				log.Printf("applied courier %d to order %s", courier.ID, order.ID)
			}
		}
	}
}

func createAssignmentInTransaction(
	ctx context.Context,
	txRunner txRunner,
	courierRepository courierRepository,
	deliveryRepository deliveryRepository,
	factory deliveryCalculatorFactory,
	orderID string,
) (model.Courier, model.Delivery, error) {
	var (
		courier  model.Courier
		delivery model.Delivery
	)

	err := txRunner.Run(ctx, func(txCtx context.Context) error {
		c, err := courierRepository.FindAvailableCourier(txCtx)
		if err != nil {
			if errors.Is(err, courierrepoerrors.ErrCouriersBusy) {
				return ErrCouriersBusy
			}
			return err
		}

		dc := factory.GetDeliveryCalculator(c.TransportType)
		if dc == nil {
			return ErrUnknownTransportType
		}
		deliveryDomain := model.Delivery{
			OrderID:    orderID,
			CourierID:  c.ID,
			AssignedAt: time.Now(),
			Deadline:   dc.CalculateDeadline(),
		}

		d, err := deliveryRepository.CreateDelivery(txCtx, deliveryDomain)
		if err != nil {
			if errors.Is(err, deliveryrepoerrors.ErrOrderIDExists) {
				return ErrOrderIDExists
			}
			return err
		}

		c.ChangeStatus(model.CourierStatusBusy)
		if err := courierRepository.UpdateCourier(txCtx, c); err != nil {
			return err
		}

		courier = c
		delivery = d
		return nil
	})
	if err != nil {
		return model.Courier{}, model.Delivery{}, err
	}

	return courier, delivery, nil
}
