package publisher

import (
	"log"

	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderActionPublisher struct {
	ch *amqp.Channel
}

func NewOrderActionPublisher(ch *amqp.Channel) *OrderActionPublisher {
	return &OrderActionPublisher{ch: ch}
}

func (o *OrderActionPublisher) PublishOrderAction(exchangeName, routingKey string, order amqpmodel.OrderAction) error {
	err := publisher.PublishMessage(o.ch, exchangeName, routingKey, order)
	if err != nil {
		log.Printf("Failed to publish order action: %v", err)
		return err
	}
	return nil
}
