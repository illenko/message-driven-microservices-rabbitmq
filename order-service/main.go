package main

import (
	"log"
	"net/http"

	"github.com/illenko/order-service/internal/consumer"
	"github.com/illenko/order-service/internal/handler"
	"github.com/illenko/order-service/internal/mapper"
	"github.com/illenko/order-service/internal/messaging"
	"github.com/illenko/order-service/internal/publisher"
	"github.com/illenko/order-service/internal/repository"
	"github.com/illenko/order-service/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	messagingConfig := messaging.Config{
		ConnectionUrl:       "amqp://user:password@localhost:5672/",
		OrderActionExchange: "order-action-exchange",
		OrderResultExchange: "order-result-exchange",
		DLXExchange:         "dlx-exchange",
		ProductReservationQueue: messaging.QueueConfig{
			Name: "product-reservation-action-queue",
			TTL:  15000,
			DLX:  "dlx-exchange",
		},
		ProductCancellationQueue: messaging.QueueConfig{
			Name: "product-cancellation-action-queue",
		},
		PaymentQueue: messaging.QueueConfig{
			Name: "payment-action-queue",
			TTL:  60000,
			DLX:  "dlx-exchange",
		},
		ProductReservationResultQueue: messaging.QueueConfig{
			Name: "product-reservation-result-queue",
		},
		PaymentResultQueue: messaging.QueueConfig{
			Name: "payment-result-queue",
		},
		DLXProductReservationQueue: messaging.QueueConfig{
			Name: "dlx-product-reservation-action-queue",
		},
		DLXPaymentQueue: messaging.QueueConfig{
			Name: "dlx-payment-action-queue",
		},
	}

	conn, ch := setupMessaging(messagingConfig)
	defer conn.Close()
	defer ch.Close()

	repo := repository.NewInMemoryOrderRepository()
	orderMapper := mapper.NewOrderMapper()
	amqpMapper := mapper.NewAmqpMapper()
	orderActionPublisher := publisher.NewOrderActionPublisher(ch)

	orderService := service.NewOrderService(messagingConfig, orderActionPublisher, repo, orderMapper, amqpMapper)

	startConsumers(ch, messagingConfig, orderService)
	startHttpServer(orderService)
}

func setupMessaging(config messaging.Config) (*amqp.Connection, *amqp.Channel) {
	broker := messaging.NewBroker(config)

	conn, ch, err := broker.Setup()
	if err != nil {
		log.Fatalf("Failed to setup RabbitMQ: %v", err)
	}

	if err := broker.SetupDLX(ch); err != nil {
		log.Fatalf("Failed to setup DLX: %v", err)
	}

	if err := broker.SetupExchangesAndQueues(ch); err != nil {
		log.Fatalf("Failed to setup exchanges and queues: %v", err)
	}

	return conn, ch
}

func startConsumers(ch *amqp.Channel, messagingConfig messaging.Config, orderService *service.OrderService) {
	orderActionResultConsumer := consumer.NewOrderActionResultConsumerImpl(ch, messagingConfig, orderService)
	expiredActionConsumer := consumer.NewExpiredActionConsumerImpl(ch, messagingConfig, orderService)

	go orderActionResultConsumer.ConsumePaymentMessages()
	go orderActionResultConsumer.ConsumeProductReservationMessages()
	go expiredActionConsumer.ConsumeExpiredProductReservation()
	go expiredActionConsumer.ConsumeExpiredPayment()
}

func startHttpServer(orderService *service.OrderService) {
	orderHandler := handler.NewOrderHandlerImpl(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrder)
	_ = http.ListenAndServe(":8080", mux)
}
