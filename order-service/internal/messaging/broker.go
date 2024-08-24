package messaging

import (
	"github.com/illenko/common/connection"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	config Config
}

func NewBroker(config Config) *Broker {
	return &Broker{config: config}
}

func (b *Broker) Setup() (*amqp.Connection, *amqp.Channel, error) {
	conn, ch, err := connection.ConnectToRabbitMQ(b.config.ConnectionUrl)
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, nil
}

func (b *Broker) SetupDLX(ch *amqp.Channel) error {
	if err := b.declareExchange(ch, b.config.DLXExchange, "direct"); err != nil {
		return err
	}

	return b.setupQueues(ch, []QueueConfig{
		b.config.DLXProductReservationQueue,
		b.config.DLXPaymentQueue,
	})
}

func (b *Broker) SetupExchangesAndQueues(ch *amqp.Channel) error {
	if err := b.setupExchangeAndQueues(ch, b.config.OrderActionExchange, []QueueConfig{
		b.config.ProductReservationQueue,
		b.config.ProductCancellationQueue,
		b.config.PaymentQueue,
	}); err != nil {
		return err
	}

	return b.setupExchangeAndQueues(ch, b.config.OrderResultExchange, []QueueConfig{
		b.config.ProductReservationResultQueue,
		b.config.PaymentResultQueue,
	})
}

func (b *Broker) declareExchange(ch *amqp.Channel, exchangeName, exchangeType string) error {
	return ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (b *Broker) setupQueues(ch *amqp.Channel, queues []QueueConfig) error {
	for _, queueConfig := range queues {
		_, err := ch.QueueDeclare(
			queueConfig.Name,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		err = ch.QueueBind(
			queueConfig.Name,
			queueConfig.Name,
			b.config.DLXExchange,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) setupExchangeAndQueues(ch *amqp.Channel, exchangeName string, queues []QueueConfig) error {
	err := b.declareExchange(ch, exchangeName, "direct")
	if err != nil {
		return err
	}

	for _, queueConfig := range queues {
		args := amqp.Table{}
		if queueConfig.TTL > 0 {
			args["x-message-ttl"] = int32(queueConfig.TTL)
		}
		if queueConfig.DLX != "" {
			args["x-dead-letter-exchange"] = queueConfig.DLX
			args["x-dead-letter-routing-key"] = "dlx-" + queueConfig.Name
		}

		_, err = ch.QueueDeclare(
			queueConfig.Name,
			true,
			false,
			false,
			false,
			args,
		)
		if err != nil {
			return err
		}

		err = ch.QueueBind(
			queueConfig.Name,
			queueConfig.Name,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
