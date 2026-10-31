package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anshggss/dhoondo/internal/input"
	"github.com/anshggss/dhoondo/internal/pubsub"
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

	_, err = pubsub.DeclareAndBind(serverChan, "dhoondo_websites", "durable")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println("Dhoondo is up")
	_ = input.TakeInput()
}
