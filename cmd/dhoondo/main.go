package main

import (
	"log"
	"os"

	"github.com/anshggss/dhoondo/internal/input"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	rabbitserver := os.Getenv("RABBIT_MQ_SERVER")
	conn, err := amqp.Dial(rabbitserver)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	serverChan, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = pubsub.DeclareAndBind(serverChan, routing.DhoondoQueueName, routing.DhoondoWebsitesExchange, routing.DhoondoWebsitesKey, "durable")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	websites := input.TakeInput()
	// TODO - Create a helper function that sends the input websites to the websites queue

	err = pubsub.PublishWebsites(serverChan, routing.DhoondoWebsitesExchange, routing.DhoondoWebsitesKey, websites)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
