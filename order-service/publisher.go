package main

import (
	"encoding/json"
	"log"

	"github.com/illenko/common/amqpmodel"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishOrderAction(ch *amqp.Channel, exchangeName, routingKey string, order amqpmodel.OrderAction) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err == nil {
		log.Printf("Order action published: %v in %v with %v rk", order, exchangeName, routingKey)
	}
	return err
}
