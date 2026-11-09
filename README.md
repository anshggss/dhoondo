# Dhoondo

Dhoondo is a website indexer CLI. You give it one or more website URLs, it publishes those to RabbitMQ and a worker then goes out and crawls each page, extracts the words from it and stores a reverse index in PostgreSQL.

## How it works

- The CLI takes URLs from command line arguments.
- The main process publishes valid URLs to the dhoondo_websites RabbitMQ queue.
- A worker consumes the queue and pulls each page.
- The worker pulls the words out and de-duplicates them, then stores the websites, the words and the word-to-website links in PostgreSQL.
- The database reverse index is the basis for searching the indexed pages. There's no search CLI/API at present.

Recursive crawling is supported in the crawler helpers, though newly discovered links are not at present published back out to Rabbitmq. The worker will index the page that's been submitted by the client.

> **Warning:** If recursive crawling is ever enabled in the future then sites like Wikipedia can contain a very large number of links. Set MAX_DEPTH carefully to avoid an unexpectedly large crawl.

## Requirements

- Go 1.27.3 or compatible
- Docker and Docker Compose
- RabbitMQ
- PostgreSQL

## Configuration

Both executables will load a .env file from the working directory. Define:

```env
RABBIT_MQ_SERVER=amqp://guest:guest@localhost:5672/
POSTGRES_CONN=postgres://postgres:postgres@localhost:5432/index
MAX_DEPTH=1
USER_AGENT=dhoondo/1.0
```

Keep your credentials and any other environment specific values out of version control. .env.local can be used for local overrides, the app currently calls godotenv.Load() so it's expecting the active file to be called .env.

## Running locally

1. Start PostgreSQL:

   ```bash
   cd cmd/workers
   docker compose up -d db
   cd ../..
   ```

2. Start RabbitMQ using the helper script:

   ```bash
   ./rabbit.sh start
   ```

3. Apply the SQL files in internal/sql/schema/ to the configured PostgreSQL database.

4. Start the worker in one terminal:

   ```bash
   go run ./cmd/workers
   ```

5. Submit one or more URLs from another terminal:

   ```bash
   go run ./cmd/dhoondo https://example.com https://www.example.org
   ```

The worker logs the number of words found and the URL indexed. Stop the local services with ./rabbit.sh stop and docker compose -f cmd/workers/compose.yaml down.

## Project layout

```text
cmd/dhoondo/             URL input and RabbitMQ publisher
cmd/workers/             RabbitMQ worker entrypoint and crawler
internal/input/          URL validation and argument parsing
internal/pubsub/         RabbitMQ declaration and publishing
internal/sql/schema/     PostgreSQL schema migrations
internal/sql/queries/    SQL queries used by sqlc
internal/database/       Generated sqlc PostgreSQL code
```

## Database

The schema contains:

- `websites`: indexed URLs
- `words`: unique words extracted from pages
- `reverseindex`: relationships between words and websites

`sqlc.yaml` produces the PostgreSQL data access code in internal/database/ from the SQL files in internal/sql/queries/.

## Development

Run the Go test suite with:

```bash
go test ./...
```
