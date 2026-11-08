package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/anshggss/dhoondo/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishWebsites(q *amqp.Channel, exchange, key string, websites []models.Website) error {
	if len(websites) == 0 {
		return nil
	}
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
		fmt.Printf("Published website %s\n", website.Link)
	}
	return nil
}
