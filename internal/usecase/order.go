package usecase

import (
	"log"
	"time"
	"context"
)

type OrderUsecase struct {
	orderGateway orderGateway
	courierRepository сourierRepository
	txRunner txRunner
	deliveryRepository deliveryRepository
	deliveryCalculatorFactory deliveryCalculatorFactory
}

func NewOrderUsecase(
	orderGateway orderGateway, 
	courierRepository сourierRepository, 
	deliveryRepository deliveryRepository, 
	txRunner txRunner,
	deliveryCalculatorFactory deliveryCalculatorFactory,
) *OrderUsecase {
	return &OrderUsecase{
		orderGateway: orderGateway,
		courierRepository: courierRepository,
		deliveryRepository: deliveryRepository,
		txRunner: txRunner,
		deliveryCalculatorFactory: deliveryCalculatorFactory,
	}
}

func (u *OrderUsecase) ProcessOrders(ctx context.Context, interval time.Duration) {
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