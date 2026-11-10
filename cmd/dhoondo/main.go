package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/input"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/jackc/pgx/v5"
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

	args := os.Args

	if len(args) < 1 {
		fmt.Println("give something ra")
	}
	switch args[1] {

	case "index":
		if len(args) < 2 {
			fmt.Println("Wow goodjob")
			return
		}
		websites := input.WebsiteHandler(args[2:])
		err = pubsub.PublishWebsites(serverChan, routing.DhoondoWebsitesExchange, routing.DhoondoWebsitesKey, websites)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	case "query":
		connstr := os.Getenv("POSTGRES_CONN")
		dbconn, err := pgx.Connect(context.Background(), connstr)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		qry := database.New(dbconn)
		matchedWebsites := input.QueryHandler(args[2:], qry)
		if matchedWebsites == nil {
			fmt.Println("Nothing matched")
			return
		}
		for _, webs := range matchedWebsites {
			fmt.Println(webs)
		}
	default:
		input.PrintHelp()
	}

}
