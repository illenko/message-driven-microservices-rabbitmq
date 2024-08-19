package main

import (
	"log"

	"github.com/illenko/common/connection"
)

func main() {
	conn, ch, err := connection.ConnectToRabbitMQ("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	go consumePaymentMessages(ch)
	select {}
}
