package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	re "courier-service/internal/gateway/retry"
	"courier-service/internal/model"
	l "courier-service/pkg/logger"
	pb "courier-service/proto/order"
)

type Gateway struct {
	client    client
	retryexec retryexec
	logger    l.LoggerInterface
}

func NewGateway(client client, rexec retryexec, l l.LoggerInterface) *Gateway {
	return &Gateway{
		client:    client,
		retryexec: rexec,
		logger:    l,
	}
}

func (g *Gateway) GetOrders(ctx context.Context, from time.Time) ([]model.Order, error) {
	var orders *pb.GetOrdersResponse

	err := g.retryexec.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := g.client.GetOrders(ctx, &pb.GetOrdersRequest{
			From: timestamppb.New(from),
		})
		if err != nil {
			return fmt.Errorf("failed to get orders: %w", err)
		}
		if resp == nil {
			return errors.New("no orders found")
		}
		orders = resp
		return nil
	})

	if err != nil {
		return nil, err
	}

	ordersList := make([]model.Order, 0, len(orders.Orders))
	for _, order := range orders.Orders {
		ordersList = append(ordersList, orderModelFromProto(order))
	}
	return ordersList, nil
}

func (g *Gateway) GetOrderById(ctx context.Context, id string) (model.Order, error) {
	var order *pb.GetOrderByIdResponse

	g.logger.Infof("sending grpc request for order id %s", id)
	err := g.retryexec.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := g.client.GetOrderById(ctx, &pb.GetOrderByIdRequest{
			Id: id,
		})
		if err != nil {
			return fmt.Errorf("failed to get order by id: %w", err)
		}
		if resp == nil {
			return errors.New("order not found")
		}
		order = resp
		return nil
	})

	if err != nil {
		if errors.Is(err, re.ErrMaxAttemptsExceeded) {
			return model.Order{}, ErrRetryLimitExceeded
		}
		return model.Order{}, err
	}

	return orderModelFromProto(order.Order), nil
}
