package order

import (
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"

	"courier-service/internal/model"
	changed "courier-service/internal/usecase/order/changed"
	logger "courier-service/pkg/logger"
)

type OrderStatusChangedHandler struct {
	useCase orderChangedUseCase
	logger  logger.LoggerInterface
}

func NewOrderStatusChangedHandler(useCase orderChangedUseCase, logger logger.LoggerInterface) *OrderStatusChangedHandler {
	return &OrderStatusChangedHandler{useCase: useCase, logger: logger}
}

func (h *OrderStatusChangedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *OrderStatusChangedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *OrderStatusChangedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := session.Context()

		var msg orderChangedDto
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			session.MarkMessage(message, "")
			continue
		}

		status := model.OrderStatus(msg.Status)

		h.logger.Infof("fetched order with id %s and status %s", msg.OrderID, status)

		if err := h.useCase.HandleOrderStatusChanged(ctx, status, msg.OrderID); err != nil {
			if errors.Is(err, changed.ErrOrderStatusMismatch) {
				session.MarkMessage(message, "")
				continue
			}
			h.logger.Errorf("order.changed handler: failed to process order: %v", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}
