package publisher

import (
	"encoding/json"
	"log"

	"github.com/illenko/common/amqpmodel"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishOrderResult(ch *amqp.Channel, exchangeName, routingKey string, orderResult amqpmodel.OrderActionResult) error {
	body, err := json.Marshal(orderResult)
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
		log.Printf("Order result published: %v in %v with %v rk", orderResult, exchangeName, routingKey)
	}
	return err
}
