package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/models"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/html"
)

func AddListener(qchan *amqp.Channel, qName string) error {
	qchan.Qos(10, 0, true)
	del, err := qchan.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = godotenv.Load()
	if err != nil {
		return err
	}
	psql := os.Getenv("POSTGRES_CONN")
	dbconn, err := pgx.Connect(context.Background(), psql)
	if err != nil {
		return err
	}
	queries := database.New(dbconn)
	for d := range del {
		var val models.Website
		if err := json.Unmarshal(d.Body, &val); err != nil {
			// Not qing cause we have shit ton of things
			d.Nack(false, true)
			continue
		}
		cnt, err := queries.FindWebsite(context.Background(), pgtype.Text{
			String: val.Link,
			Valid:  true,
		})
		if err != nil {
			d.Nack(false, true)
			continue
		}
		if cnt > 0 {
			d.Nack(false, false)
			continue
		}
		// If there is any row present, we do not crawl and instead return

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
			d.Nack(false, true)
			continue
		}
		defer resp.Body.Close()
		// Before parsing, now we insert the website in the database
		root, err := html.Parse(resp.Body)
		if err != nil {
			d.Nack(false, false)
			continue
		}
		links := ExtractLinks(root, 0)
		if links == nil {
			d.Nack(false, false)
			continue
		}
		websites := []models.Website{}
		fmt.Println("Printing websites")
		for _, v := range links {
			base, err := url.Parse(val.Link)
			if err != nil {
				d.Nack(false, true)
				continue
			}
			href, err := url.Parse(v)
			if err != nil {
				d.Nack(false, true)
				continue
			}
			resolved := base.ResolveReference(href)
			curLink := models.Website{
				Link: resolved.String(),
			}
			websites = append(websites, curLink)
			fmt.Printf("Added %s to publishing queue\n", curLink.Link)
		}
		err = pubsub.PublishWebsites(qchan,
			routing.DhoondoWebsitesExchange,
			routing.DhoondoWebsitesKey,
			websites)
		if err != nil {
			fmt.Println(err.Error())
			d.Nack(false, true)
			continue
		}
		id, err := uuid.NewV7()
		if err != nil {
			fmt.Println(err.Error())
			d.Nack(false, true)
			continue
		}
		arg := database.InsertWebsiteParams{
			ID: pgtype.UUID{
				Bytes: id,
				Valid: true,
			},
			Url: pgtype.Text{
				String: val.Link,
				Valid:  true,
			},
		}
		err = queries.InsertWebsite(context.Background(), arg)
		if err != nil {
			fmt.Println(err.Error())
			d.Nack(false, true)
			continue
		}
		d.Ack(false)
		fmt.Printf("Added %s to database\n", val.Link)
	}
	return nil
}
