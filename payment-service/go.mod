module github.com/illenko/payment-service

go 1.22.3

require (
	github.com/illenko/common v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

require github.com/google/uuid v1.6.0 // indirect

replace github.com/illenko/common => ../common
