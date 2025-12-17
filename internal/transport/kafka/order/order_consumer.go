package order

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"courier-service/internal/model"
)

type OrderStatusChangedHandler struct {
	factory *OrderChangedFactory
}

func NewOrderStatusChangedHandler(
	factory *OrderChangedFactory,
) *OrderStatusChangedHandler {
	return &OrderStatusChangedHandler{
		factory: factory,
	}
}

func (h *OrderStatusChangedHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (h *OrderStatusChangedHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *OrderStatusChangedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := session.Context()
		log.Printf("order.changed handler: received message: key=%s, value=%s, partition=%d, offset=%d",
			string(message.Key), string(message.Value), message.Partition, message.Offset)

		var msg orderChangedDto
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("order.changed handler: failed to unmarshal message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		msgStatus := model.OrderStatus(msg.Status)

		if msgStatus != model.OrderStatusCreated && msgStatus != model.OrderStatusCompleted && msgStatus != model.OrderStatusCancelled {
			log.Printf("order.changed handler: unknown order status: %v", msg.Status)
			session.MarkMessage(message, "")
			continue
		}

		handler := h.factory.GetOrderChanged(model.OrderStatus(msg.Status))
		if handler == nil {
			log.Printf("order.changed handler: no handler found for status: %s", msg.Status)
			session.MarkMessage(message, "")
			continue
		}

		err := handler.HandleOrderStatusChanged(ctx, model.OrderStatus(msg.Status), msg.OrderID)
		if err != nil {
			log.Printf("order.changed handler: failed to process order: %v", err)
		}

		session.MarkMessage(message, "")
	}

	return nil
}
