package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/illenko/common/mapper"

	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	"github.com/illenko/common/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

var failedItemId = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func consumeProductReservationMessages(ch *amqp.Channel) {
	consumer.ConsumeMessages(ch, "product-reservation-action-queue", func(orderAction amqpmodel.OrderAction) {
		status := amqpmodel.OrderActionResultStatusSuccess

		if orderAction.ItemID == failedItemId {
			status = amqpmodel.OrderActionResultStatusFailed
		}

		orderResult := mapper.ToOrderActionResult(orderAction, status)

		err := publisher.PublishMessage(ch, "order-result-exchange", "product-reservation-result-queue", orderResult)
		if err != nil {
			log.Printf("Failed to publish order result: %v", err)
		}

	})
}
