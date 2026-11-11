# Dhoondo

Dhoondo is a website indexer CLI. You give it one or more website URLs, it publishes those to RabbitMQ and a worker then goes out and crawls each page, extracts the words from it and stores a reverse index in PostgreSQL. You can later on look for words in the websites using index.

## How it works

- The CLI takes URLs from command line arguments.
- The main process publishes valid URLs to the dhoondo_websites RabbitMQ queue.
- A worker consumes the queue and pulls each page.
- The worker pulls the words out and de-duplicates them, then stores the websites, the words and the word-to-website links in PostgreSQL.
- The database reverse index is the basis for searching the indexed pages.

The CLI supports indexing URLs and querying the reverse index. Recursive crawling helpers are present, though newly discovered links are not currently published back to RabbitMQ. The worker indexes the page submitted by the client.

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

Keep your credentials and other environment-specific values out of version control. The application currently calls `godotenv.Load()`, so the active file must be named `.env` in the working directory.

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
   go run ./cmd/dhoondo index https://example.com https://www.example.org
   ```

The worker logs the number of words found and the URL indexed. Stop the local services with:

```bash
./rabbit.sh stop
docker compose -f cmd/workers/compose.yaml down
```

## Querying indexed websites

After a URL has been indexed, query one or more words with:

```bash
go run ./cmd/dhoondo query dhoondo search
```

The command prints the indexed URLs matching each query word. PostgreSQL and RabbitMQ must be running because the CLI connects to RabbitMQ before handling commands.

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
