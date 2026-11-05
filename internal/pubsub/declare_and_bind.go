package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(serverChan *amqp.Channel, queuename, exchange, key, durable string) (*amqp.Queue, error) {

	// Declare the exchange
	err := serverChan.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	// Declare the queue
	q, err := serverChan.QueueDeclare(queuename,
		durable == "durable",
		false, false, false, nil)
	if err != nil {
		return nil, err
	}

	// Bind the exchange to the queue with the given
	err = serverChan.QueueBind(queuename, key, exchange, true, nil)
	if err != nil {
		return nil, err
	}

	return &q, err
}
