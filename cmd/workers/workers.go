// These guys here will be chilling and waiting for work to arrive
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anshggss/dhoondo/cmd/workers/crawler"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// This function here talks to the rabbit mq server.
// For each message, it calls the crawl function and basically spawns a new crawler
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
	defer conn.Close()
	qchan, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	err = crawler.AddListener(qchan, routing.DhoondoQueueName)
	if err != nil {
		fmt.Println(err.Error())
	}
}
