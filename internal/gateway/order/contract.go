package order

import (
	"context"
	pb "courier-service/proto/order"
	"time"

	"google.golang.org/grpc"
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
