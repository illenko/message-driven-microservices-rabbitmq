package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/illenko/common/amqpmodel"
	amqp "github.com/rabbitmq/amqp091-go"
)

func consumeOrderResult(ch *amqp.Channel, queueName string, processFunc func(*amqp.Channel, amqpmodel.OrderActionResult)) {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var orderAction amqpmodel.OrderActionResult
		err := json.Unmarshal(msg.Body, &orderAction)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Received a message from %v: %v", queueName, orderAction)
		processFunc(ch, orderAction)
	}
}

func consumeProductReservationMessages(ch *amqp.Channel) {
	consumeOrderResult(ch, "product-reservation-result-queue", processProductReservationResult)
}

func processProductReservationResult(ch *amqp.Channel, orderResult amqpmodel.OrderActionResult) {
	log.Printf("Processing product reservation result: %v and other work will be done...", orderResult)
	err := publishOrderAction(ch, "order-action-exchange", "payment-action-queue", amqpmodel.OrderAction{
		ID:         orderResult.ID,
		CustomerID: uuid.New(),
		CardID:     uuid.New(),
		ItemID:     uuid.New(),
		Price:      100,
	})
	if err != nil {
		log.Fatalf("Failed to publish order action: %v", err)
		return
	}
}

func consumePaymentMessages(ch *amqp.Channel) {
	consumeOrderResult(ch, "payment-result-queue", processPaymentResult)
}

func processPaymentResult(_ *amqp.Channel, orderResult amqpmodel.OrderActionResult) {
	log.Printf("Processing payment result: %v and other work will be done...", orderResult)
}
