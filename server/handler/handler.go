package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anshggss/dhoondo/internal/database"
	"github.com/anshggss/dhoondo/internal/input"
	"github.com/anshggss/dhoondo/internal/pubsub"
	"github.com/anshggss/dhoondo/internal/routing"
	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"
)

type publishRequest struct {
	Websites []string `json:"websites"`
}

type searchRequest struct {
	Words []string `json:"words"`
}

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func Publish(channel *amqp.Channel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var request publishRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		websites := input.WebsiteHandler(request.Websites)
		if len(websites) == 0 {
			http.Error(w, "websites must contain at least one valid URL", http.StatusBadRequest)
			return
		}

		if err := pubsub.PublishWebsites(channel, routing.DhoondoWebsitesExchange, routing.DhoondoWebsitesKey, websites); err != nil {
			http.Error(w, "failed to publish websites", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, map[string]int{"published": len(websites)})
	}
}

func Search(dbq *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var request searchRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if len(request.Words) == 0 {
			http.Error(w, "words must contain at least one search term", http.StatusBadRequest)
			return
		}

		websites := make([]string, 0)
		for _, word := range request.Words {
			matches, err := dbq.QueryIndex(r.Context(), pgtype.Text{
				String: word,
				Valid:  true,
			})
			if err != nil {
				http.Error(w, "failed to search websites", http.StatusInternalServerError)
				return
			}
			for _, match := range matches {
				if match.Valid {
					websites = append(websites, match.String)
				}
			}
		}

		writeJSON(w, http.StatusOK, map[string][]string{"websites": websites})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func setCORSHeaders(w http.ResponseWriter, methods string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", methods)
}
