package ordermonitoring

import (
	"context"
	"log"
	"time"
)

type OrderMonitoringUseCase struct {
	orderGateway              orderGateway
	courierRepository         courierRepository
	txRunner                  txRunner
	deliveryRepository        deliveryRepository
	deliveryCalculatorFactory deliveryCalculatorFactory
	assignUseCase             assignUseCase
}

func NewOrderMonitoringUseCase(
	orderGateway orderGateway,
	courierRepository courierRepository,
	deliveryRepository deliveryRepository,
	txRunner txRunner,
	deliveryCalculatorFactory deliveryCalculatorFactory,
	assignUseCase assignUseCase,
) *OrderMonitoringUseCase {
	return &OrderMonitoringUseCase{
		orderGateway:              orderGateway,
		courierRepository:         courierRepository,
		deliveryRepository:        deliveryRepository,
		txRunner:                  txRunner,
		deliveryCalculatorFactory: deliveryCalculatorFactory,
		assignUseCase:             assignUseCase,
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
				assignment, err := u.assignUseCase.Assign(ctx, order.ID)
				if err != nil {
					log.Printf("failed to create assignment for order %s: %v\n", order.ID, err)
					continue
				}
				log.Printf("applied courier %d to order %s", assignment.CourierID, order.ID)
			}
		}
	}
}
