package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/connection"
)

func main() {
	conn, ch, err := connection.ConnectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	err = declareDLXExchange(ch)
	if err != nil {
		log.Fatalf("Failed to declare DLX exchange: %v", err)
	}

	err = setupDLXQueues(ch, []QueueConfig{
		{Name: "dlx-product-reservation-action-queue"},
		{Name: "dlx-payment-action-queue"},
	})
	if err != nil {
		log.Fatalf("Failed to setup DLX queues: %v", err)
	}

	err = setupExchangeAndQueues(ch, "order-action-exchange", []QueueConfig{
		{Name: "product-reservation-action-queue", TTL: 15000, DLX: "dlx-exchange"},
		{Name: "product-cancellation-action-queue"},
		{Name: "payment-action-queue", TTL: 60000, DLX: "dlx-exchange"},
	})
	if err != nil {
		log.Fatalf("Failed to setup order-action-exchange: %v", err)
	}

	err = setupExchangeAndQueues(ch, "order-result-exchange", []QueueConfig{
		{Name: "product-reservation-result-queue"},
		{Name: "payment-result-queue"},
	})
	if err != nil {
		log.Fatalf("Failed to setup order-result-exchange: %v", err)
	}

	orderAction := amqpmodel.OrderAction{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		CardID:     uuid.New(),
		ItemID:     uuid.New(),
		Price:      100,
	}

	err = publishOrderAction(ch, "order-action-exchange", "product-reservation-action-queue", orderAction)
	if err != nil {
		log.Fatalf("Failed to publish order: %v", err)
	}

	go consumeExpiredPayment(ch)
	go consumeExpiredProductReservation(ch)
	go consumeProductReservationMessages(ch)
	go consumePaymentMessages(ch)

	select {}
}
