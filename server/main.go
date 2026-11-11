package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/anshggss/dhoondo/server/handler"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dbConn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_CONN"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(context.Background())
	dbq := database.New(dbConn)

	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	if _, err := pubsub.DeclareAndBind(
		channel,
		routing.DhoondoQueueName,
		routing.DhoondoWebsitesExchange,
		routing.DhoondoWebsitesKey,
		"durable",
	); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.Health)
	mux.HandleFunc("POST /publish", handler.Publish(channel))
	mux.HandleFunc("OPTIONS /publish", handler.Publish(channel))
	mux.HandleFunc("POST /search/", handler.Search(dbq))
	mux.HandleFunc("OPTIONS /search/", handler.Search(dbq))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
