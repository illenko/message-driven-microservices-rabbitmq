package main

import (
	"log"

	"github.com/illenko/common/mapper"

	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	"github.com/illenko/common/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

func consumePaymentMessages(ch *amqp.Channel) {
	consumer.ConsumeOrderAction(ch, "payment-action-queue", processPayment)
}

func processPayment(ch *amqp.Channel, orderAction amqpmodel.OrderAction) {
	orderResult := mapper.ToOrderActionResult(orderAction, amqpmodel.OrderActionResultStatusSuccess)

	err := publisher.PublishOrderResult(ch, "order-result-exchange", "payment-result-queue", orderResult)
	if err != nil {
		log.Printf("Failed to publish order result: %v", err)
	}
}
