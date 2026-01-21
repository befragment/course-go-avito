package order

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "courier-service/proto/order"
)

type client interface {
	GetOrders(ctx context.Context, in *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error)
	GetOrderById(ctx context.Context, in *pb.GetOrderByIdRequest, opts ...grpc.CallOption) (*pb.GetOrderByIdResponse, error)
}

type retryexec interface {
	Execute(fn func() error) error
	ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error
	ExecuteWithCallback(
		fn func() error,
		onRetry func(attempt int, err error, delay time.Duration),
	) error
}
