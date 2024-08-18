package main

import (
	"log"

	"github.com/illenko/common/connection"
)

func main() {
	conn, ch, err := connection.ConnectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	go consumePaymentMessages(ch)
	select {}
}
