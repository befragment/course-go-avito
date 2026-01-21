package order

import (
	"courier-service/internal/model"
	pb "courier-service/proto/order"
)

func orderModelFromProto(o *pb.Order) model.Order {
	return model.Order{
		ID:        o.Id,
		Status:    model.OrderStatus(o.Status),
		CreatedAt: o.CreatedAt.AsTime(),
	}
}
