package consumer

import (
	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	"github.com/illenko/order-service/internal/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ActionResultProcessor interface {
	ProcessProductReservationResult(orderResult amqpmodel.OrderActionResult)
	ProcessPaymentResult(orderResult amqpmodel.OrderActionResult)
}

type OrderActionResultConsumer interface {
	ConsumeProductReservationMessages()
	ConsumePaymentMessages()
}

type OrderActionResultConsumerImpl struct {
	ch              *amqp.Channel
	messagingConfig messaging.Config
	processor       ActionResultProcessor
}

func NewOrderActionResultConsumerImpl(ch *amqp.Channel, messagingConfig messaging.Config, processor ActionResultProcessor) *OrderActionResultConsumerImpl {
	return &OrderActionResultConsumerImpl{ch: ch, messagingConfig: messagingConfig, processor: processor}
}

func (o *OrderActionResultConsumerImpl) ConsumeProductReservationMessages() {
	consumer.ConsumeMessages(o.ch, o.messagingConfig.ProductReservationResultQueue.Name, o.processor.ProcessProductReservationResult)
}

func (o *OrderActionResultConsumerImpl) ConsumePaymentMessages() {
	consumer.ConsumeMessages(o.ch, o.messagingConfig.PaymentResultQueue.Name, o.processor.ProcessPaymentResult)
}
