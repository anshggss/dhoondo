package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(serverChan *amqp.Channel, queuename, durable string) (*amqp.Queue, error) {

	q, err := serverChan.QueueDeclare(queuename,
		durable == "durable",
		false, false, false, nil)
	return &q, err
}
