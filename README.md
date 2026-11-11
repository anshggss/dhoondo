# Dhoondo

Dhoondo indexes websites through a RabbitMQ-backed crawler and lets you
search the collected words through a web interface.

## Requirements

- Go 1.27.3 or compatible
- Bun
- Docker

## Configuration

Create a `.env` file in the repository root:

```env
RABBIT_MQ_SERVER=amqp://guest:guest@localhost:5672/
POSTGRES_CONN=postgresql://postgres:postgres@localhost:5432/index
```

The Go server loads this file when it starts. `PORT` is optional and defaults
to `8080`.

## Start the dependencies

Start PostgreSQL:

```bash
docker compose -f cmd/workers/compose.yaml up -d db
```

Start RabbitMQ:

```bash
./rabbit.sh start
```

## Start the server

From the repository root:

```bash
go run ./server
```

The server listens on `http://localhost:8080`.

## Start the crawler worker

In another terminal, from the repository root:

```bash
go run ./cmd/workers
```

The worker listens to RabbitMQ for published websites, crawls each website,
and stores its words and reverse index in PostgreSQL. Keep the worker running
while using the frontend or `/publish`; the server only places websites on the
RabbitMQ queue.

### API endpoints

Publish websites for crawling:

```http
POST /publish
Content-Type: application/json

{"websites":["https://example.com","https://developer.mozilla.org"]}
```

Search indexed websites:

```http
POST /search/
Content-Type: application/json

{"words":["architecture"]}
```

Both endpoints return JSON and support browser CORS requests.

## Start the website frontend

In another terminal:

```bash
cd app
bun install
bun dev
```

Open the URL printed by Bun. The frontend defaults to the Go server at
`http://localhost:8080`; set `BUN_PUBLIC_API_URL` if the server runs elsewhere.

The dispatch page submits website batches through `/publish`. Use the
**Search the indexed web** link to open `/search` and query indexed words.

## Checks

Run the Go tests:

```bash
go test ./...
```

Type-check the frontend:

```bash
cd app
bun run typecheck
```

## Architecture

![Architecture diagram](references/architecture/proposedarchitecture.png)
