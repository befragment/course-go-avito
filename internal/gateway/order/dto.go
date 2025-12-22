package order

import (
	pb "courier-service/proto/order"
	"courier-service/internal/model"
)

func orderModelFromProto(o *pb.Order) model.Order {
    return model.Order{
        ID: o.Id,
        Status: model.OrderStatus(o.Status),
        CreatedAt: o.CreatedAt.AsTime(),
    }
}