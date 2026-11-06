package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/anshggss/dhoondo/internal/models"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/html"
)

func AddListeners(qchan *amqp.Channel, qName string) error {
	del, err := qchan.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = godotenv.Load()
	if err != nil {
		return err
	}
	_ = os.Getenv("POSTGRES_CONN")
	for d := range del {
		go func() {
			var val models.Website
			if err := json.Unmarshal(d.Body, &val); err != nil {
				fmt.Println(err.Error())
			}
			d.Ack(false)
			// The job is picked
			// TODO ==
			// Start crawling the webpage, index it in the database
			// and add newly discovered websites to the queue
			// Now storing the we can index it but the issue would be where are we going to store it?
			// Storing it in a SQL like database wouldn't be the most gracious way of handling this problem
			// the reason being, the searching would be mostly linear.
			// If only we had a database that stores key value pairs and persists them onto
			// the disk
			// TODO FUTURE IMPROVEMENTS ==
			// Switch from POSTGRES to some key value type store or some better alternative
			// for storing the indexes
			resp, err := http.Get(val.Link)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			root, err := html.Parse(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			links := ExtractLinks(root, 0)
			if links == nil {
				fmt.Println("Found no links embedded in the text")
				return
			}
			websites := []models.Website{}
			fmt.Println("Printing websites")
			for _, v := range links {
				link := ""
				if len(v) <= 1 {
					continue
				}
				fmt.Printf("Current string length %d\n", len(v))
				if v[0] == '/' {
					link = val.Link + v
				} else {
					link = v
				}
				temp := models.Website{
					Link: link,
				}
				websites = append(websites, temp)
			}
			err = pubsub.PublishWebsites(qchan,
				routing.DhoondoWebsitesExchange,
				routing.DhoondoWebsitesKey,
				websites)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}()
	}
	return nil
}
