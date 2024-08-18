package consumer

import (
	"encoding/json"
	"log"

	"github.com/illenko/common/amqpmodel"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeOrderAction(ch *amqp.Channel, queueName string, processFunc func(*amqp.Channel, amqpmodel.OrderAction)) {
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
		var orderAction amqpmodel.OrderAction
		err := json.Unmarshal(msg.Body, &orderAction)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Received a message from %v: %v", queueName, orderAction)
		processFunc(ch, orderAction)
	}
}
