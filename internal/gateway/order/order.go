package order

import (
	"context"
	"errors"
	"time"
	"fmt"
	pb "courier-service/proto/order"
	"courier-service/internal/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gateway struct {
	client client
}

func NewGateway(client client) *Gateway {
	return &Gateway{client: client}
}

func (g *Gateway) GetOrders(ctx context.Context, from time.Time) ([]model.Order, error) {
	orders, err := g.client.GetOrders(ctx, &pb.GetOrdersRequest{
		From: timestamppb.New(from),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	if orders == nil {
		return nil, errors.New("no orders found")
	}

	ordersList := make([]model.Order, len(orders.Orders))
	for i, order := range orders.Orders {
		ordersList[i] = orderFromProto(order)
	}
	return ordersList, nil
}

func (g *Gateway) GetOrderById(ctx context.Context, id string) (model.Order, error) {
	order, err := g.client.GetOrderById(ctx, &pb.GetOrderByIdRequest{
		Id: id,
	})
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order by id: %w", err)
	}
	if order == nil {
		return model.Order{}, errors.New("order not found")
	}
	return orderFromProto(order.Order), nil
}