package order

import (
	pb "courier-service/proto/order"
	"courier-service/internal/model"
)

func orderFromProto(o *pb.Order) model.Order {
    items := make([]model.Item, len(o.Items))
    for i, it := range o.Items {
        items[i] = model.Item{
            Name:     it.Name,
            Price:    it.Price,
            Quantity: it.Quantity,
        }
    }

    return model.Order{
        ID: o.Id,
        UserID: o.UserId,
        OrderNumber: o.OrderNumber,
        FIO: o.Fio,
        RestaurantID: o.RestaurantId,
        Items: items,
        TotalPrice: o.TotalPrice,
        Address: model.Address{
            Street: o.Address.Street,
            House: o.Address.House,
            Apartment: o.Address.Apartment,
            Floor: o.Address.Floor,
            Comment: o.Address.Comment,
        },
        Status: model.OrderStatus(o.Status),
        CreatedAt: o.CreatedAt.AsTime(),
        UpdatedAt: o.UpdatedAt.AsTime(),
    }
}