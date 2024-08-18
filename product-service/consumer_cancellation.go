package main

import (
	"log"

	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func consumeProductCancellationMessages(ch *amqp.Channel) {
	consumer.ConsumeOrderAction(ch, "product-cancellation-action-queue", processProductCancellation)
}

func processProductCancellation(_ *amqp.Channel, orderAction amqpmodel.OrderAction) {
	log.Printf("Processing product cancellation: %v and other work will be done...", orderAction)
}
