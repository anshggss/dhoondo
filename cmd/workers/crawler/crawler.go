package crawler

import (
	"encoding/json"
	"fmt"

	"github.com/anshggss/dhoondo/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AddListeners(qchan *amqp.Channel, qName string) error {
	del, err := qchan.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for d := range del {
		go func() {
			var val models.Website
			if err := json.Unmarshal(d.Body, &val); err != nil {
				fmt.Println(err.Error())
			}
			d.Ack(false)
			// The job is picked
			// TODO
			// Now we need to start crawling the webpage, index it in the database
			// and add newly discovered websites to the queue
		}()
	}
	return nil
}
