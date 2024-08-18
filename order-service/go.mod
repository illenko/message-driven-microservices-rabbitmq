module github.com/illenko/order-service

go 1.22.3

require (
	github.com/google/uuid v1.6.0
	github.com/illenko/common v0.0.0-00010101000000-000000000000
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace github.com/illenko/common => ../common
