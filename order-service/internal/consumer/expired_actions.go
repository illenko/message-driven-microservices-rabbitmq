package consumer

import (
	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	"github.com/illenko/order-service/internal/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExpiredActionProcessor interface {
	ProcessExpiredProductReservation(orderAction amqpmodel.OrderAction)
	ProcessExpiredPayment(orderAction amqpmodel.OrderAction)
}

type ExpiredActionConsumer interface {
	ConsumeExpiredProductReservation()
	ConsumeExpiredPayment()
}

type ExpiredActionConsumerImpl struct {
	ch              *amqp.Channel
	messagingConfig messaging.Config
	processor       ExpiredActionProcessor
}

func NewExpiredActionConsumerImpl(ch *amqp.Channel, messagingConfig messaging.Config, processor ExpiredActionProcessor) *ExpiredActionConsumerImpl {
	return &ExpiredActionConsumerImpl{ch: ch, messagingConfig: messagingConfig, processor: processor}
}

func (e *ExpiredActionConsumerImpl) ConsumeExpiredProductReservation() {
	consumer.ConsumeMessages(e.ch, e.messagingConfig.DLXProductReservationQueue.Name, e.processor.ProcessExpiredProductReservation)
}

func (e *ExpiredActionConsumerImpl) ConsumeExpiredPayment() {
	consumer.ConsumeMessages(e.ch, e.messagingConfig.DLXPaymentQueue.Name, e.processor.ProcessExpiredPayment)
}
