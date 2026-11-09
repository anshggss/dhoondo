package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	workerhelper "github.com/anshggss/dhoondo/cmd/workers/crawler/helper"
	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/models"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
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
	mD := os.Getenv("MAX_DEPTH")
	userAgent := os.Getenv("USER_AGENT")
	maxDepth, err := strconv.Atoi(mD)
	if err != nil {
		return err
	}
	fmt.Printf("Max Depth at %d\n", maxDepth)

	dbconn, err := pgx.Connect(context.Background(), psql)
	if err != nil {
		return err
	}
	queries := database.New(dbconn)
	for d := range del {
		var val models.Website
		if err := json.Unmarshal(d.Body, &val); err != nil {
			d.Nack(false, true)
			continue
		}
		if val.Depth > maxDepth {
			d.Ack(false)
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

		// Here we need to create request and set a User-Agent header for the get request
		// because many websites need a user-agent for automated requests
		req, err := http.NewRequestWithContext(
			context.Background(),
			"GET",
			val.Link,
			nil)

		req.Header.Set(
			"User-Agent",
			userAgent)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			d.Nack(false, false)
			continue
		}

		header := resp.Header.Get("Content-Type")
		fmt.Printf("Header type is %s\n", header)
		// var root *html.Node
		// if header == "text/html" {
		root, err := html.Parse(resp.Body)
		if err != nil {
			d.Nack(false, false)
			resp.Body.Close()
			continue
		}
		// } else if header == "text/plain" {
		// 	bts, err := io.ReadAll(resp.Body)
		// 	if err != nil {
		// 		continue
		// 	}
		// 	str := string(bts)
		// 	fmt.Println(str)
		// root, err = html.Parse(strings.NewReader(resp.Body))
		// if err != nil {
		// 	d.Nack(false, false)
		// 	resp.Body.Close()
		// 	continue
		// }
		// } else {
		// 	fmt.Printf("Unknown header type %s\n", header)
		// 	d.Nack(false, false)
		// 	resp.Body.Close()
		// 	continue
		// }

		links := workerhelper.ExtractLinks(root, 0)
		if links == nil {
			d.Nack(false, false)
			resp.Body.Close()
			continue
		}
		fmt.Printf("Number of links found: %d\n", len(links))

		links = workerhelper.DeDuplicate(links)

		websites, err := workerhelper.ConvertToWebsite(links, val)
		if err != nil {
			d.Nack(false, false)
			resp.Body.Close()
			continue
		}

		err = pubsub.PublishWebsites(qchan,
			routing.DhoondoWebsitesExchange,
			routing.DhoondoWebsitesKey,
			websites)
		if err != nil {
			fmt.Println(err.Error())
			d.Nack(false, true)
			resp.Body.Close()
			continue
		}

		err = workerhelper.AddWebsiteToDatabase(val, queries)
		if err != nil {
			d.Nack(false, true)
			resp.Body.Close()
			continue
		}
		d.Ack(false)
		resp.Body.Close()
		fmt.Printf("Added %s to database\n", val.Link)
	}
	return nil
}
