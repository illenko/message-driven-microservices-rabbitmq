package main

import (
	"log"

	"github.com/illenko/common/amqpmodel"
	"github.com/illenko/common/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func consumeExpiredProductReservation(ch *amqp.Channel) {
	consumer.ConsumeOrderAction(ch, "dlx-product-reservation-action-queue", processExpiredProductReservation)
}

func processExpiredProductReservation(_ *amqp.Channel, orderAction amqpmodel.OrderAction) {
	log.Printf("Processing expired product reservation: %v and other work will be done...", orderAction)
}

func consumeExpiredPayment(ch *amqp.Channel) {
	consumer.ConsumeOrderAction(ch, "dlx-payment-action-queue", processExpiredPayment)
}

func processExpiredPayment(_ *amqp.Channel, orderAction amqpmodel.OrderAction) {
	log.Printf("Processing expired payment: %v and other work will be done...", orderAction)
}
