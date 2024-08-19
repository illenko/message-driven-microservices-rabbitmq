package publisher

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishMessage[T any](ch *amqp.Channel, exchangeName, routingKey string, message T) error {
	body, err := json.Marshal(message)
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
		log.Printf("Message published: %v in %v with %v rk", message, exchangeName, routingKey)
	}
	return err
}
