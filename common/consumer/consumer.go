package consumer

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeMessages[T any](ch *amqp.Channel, queueName string, processFunc func(T)) {
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
		var message T
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		log.Printf("Received a message from %v: %v", queueName, message)
		processFunc(message)
	}
}
