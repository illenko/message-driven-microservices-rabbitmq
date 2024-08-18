package main

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel, exchangeName, exchangeType string) error {
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

func declareDLXExchange(ch *amqp.Channel) error {
	return declareExchange(ch, "dlx-exchange", "direct")
}

func setupDLXQueues(ch *amqp.Channel, queues []QueueConfig) error {
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
			"dlx-exchange",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func setupExchangeAndQueues(ch *amqp.Channel, exchangeName string, queues []QueueConfig) error {
	err := declareExchange(ch, exchangeName, "direct")
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

type QueueConfig struct {
	Name string
	TTL  int
	DLX  string
}
