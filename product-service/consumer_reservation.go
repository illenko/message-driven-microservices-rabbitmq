package main

import (
	"github.com/illenko/common/mapper"
	"log"

	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	"github.com/illenko/common/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

func consumeProductReservationMessages(ch *amqp.Channel) {
	consumer.ConsumeOrderAction(ch, "product-reservation-action-queue", processProductReservation)
}

func processProductReservation(ch *amqp.Channel, orderAction amqpmodel.OrderAction) {
	orderResult := mapper.ToOrderActionResult(orderAction, amqpmodel.OrderActionResultStatusSuccess)

	err := publisher.PublishOrderResult(ch, "order-result-exchange", "product-reservation-result-queue", orderResult)
	if err != nil {
		log.Printf("Failed to publish order result: %v", err)
	}
}
