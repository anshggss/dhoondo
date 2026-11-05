package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/anshggss/dhoondo/internal/input"
	"github.com/anshggss/dhoondo/internal/models"
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

	_, err = pubsub.DeclareAndBind(serverChan, "dhoondo_websites", routing.DhoondoWebsitesExchange, routing.DhoondoWebsitesKey, "durable")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println("Dhoondo is up")
	websites := input.TakeInput()
	// TODO - Create a helper function that sends the input websites to the websites queue

	err = PublishWebsites(serverChan, routing.DhoondoWebsitesExchange, routing.DhoondoWebsitesKey, websites)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func PublishWebsites(q *amqp.Channel, exchange, key string, websites []models.Website) error {
	for _, website := range websites {
		val, err := json.Marshal(website)
		if err != nil {
			fmt.Printf("Error marshalling the website: %s", err.Error())
			return err
		}
		err = q.Publish(exchange, key, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        val,
		})
		if err != nil {
			return err
		}
	}
	fmt.Printf("Website %s published\n", websites[0])
	return nil
}
