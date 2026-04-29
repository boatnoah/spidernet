# Spidernet

Distributed web crawler with an API, worker, and CLI.

Spidernet accepts crawl jobs, processes pages asynchronously via a queue-backed worker, stores crawl data in Postgres, and can render a crawl graph as an image.

## Architecture

- `cmd/api`: HTTP API for submitting jobs and querying results
- `cmd/worker`: background processor that consumes queued crawl tasks
- `cmd/cli`: command-line client for API endpoints
- `postgres`: source of truth for jobs, pages, and links
- `redis`: queue backend for page crawl tasks

Flow:

1. Submit crawl job through API/CLI.
2. API persists job metadata and enqueues initial page.
3. Worker pops tasks, fetches pages, extracts links, stores results, and enqueues next depth.
4. API exposes job status and graph rendering endpoint.

## API Endpoints

- `GET /v1/health` - service health check
- `POST /v1/crawl` - submit crawl job
- `GET /v1/jobs/{jobID}/status` - fetch crawl job status
- `GET /v1/jobs/{jobID}/graph` - generate crawl graph PNG

Crawl request payload:

```json
{
  "start_url": "https://example.com",
  "depth": 2
}
```

## CLI

The CLI is an API client located at `cmd/cli`.

Top-level help:

```bash
go run ./cmd/cli --help
```

Submit a crawl:

```bash
go run ./cmd/cli crawl --depth 2 https://example.com
```

Check status:

```bash
go run ./cmd/cli status <job-id>
go run ./cmd/cli status --json <job-id>
```

Download graph:

```bash
go run ./cmd/cli graph <job-id>
go run ./cmd/cli graph --out ./crawl.png <job-id>
```

Check health:

```bash
go run ./cmd/cli health
```

Common CLI options:

- `--api-url` (default: `http://localhost:8080`)
- `--timeout` (default: `15s`)

## Local Development

### Prerequisites

- Go (matches `go.mod`)
- Docker + Docker Compose
- [`migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI

### 1) Start infrastructure

```bash
docker compose up -d
```

This starts:

- Postgres on `localhost:5432`
- Redis on `localhost:6379`
- Redis Commander on `localhost:8081`

### 2) Configure environment

This repo includes `.env` with:

```bash
DB_ADDR=postgres://admin:adminpassword@localhost/spidernet?sslmode=disable
```

You can override env vars directly when running commands (`ADDR`, `REDIS_ADDR`, etc.).

### 3) Run database migrations

```bash
make migrate-up
```

### 4) Start services

Run each process in a separate terminal.

API:

```bash
go run ./cmd/api
```

Worker:

```bash
go run ./cmd/worker
```

CLI (example):

```bash
go run ./cmd/cli crawl --depth 2 https://example.com
```
