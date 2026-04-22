# Task Scheduler

A distributed task scheduling system built in Go, using Redis for persistence and RabbitMQ as a message broker. The project follows clean architecture principles with a clear separation between domain, repository, messaging, and API layers.

-----

## Architecture Overview

```
┌─────────────┐     HTTP      ┌─────────────────────────────────┐
│   CLI Client │ ──────────►  │          REST API (Gin)          │
└─────────────┘               │  /tasks  /tasks/:id  /health     │
                               └──────────────┬──────────────────┘
                                              │ CRUD
                                    ┌─────────▼──────────┐
                                    │   Redis Repository  │
                                    │  - task:{id}        │
                                    │  - tasks (Set)      │
                                    │  - scheduled_tasks  │
                                    │    (Sorted Set)     │
                                    └─────────▲──────────┘
                                              │ poll every 30s
                               ┌──────────────┴──────────────────┐
                               │         Task Worker              │
                               │  - polls scheduled_tasks ZSet    │
                               │  - publishes to RabbitMQ         │
                               │  - consumes from tasks_queue     │
                               │  - executes via sh -c            │
                               └──────────────┬──────────────────┘
                                              │ AMQP
                                    ┌─────────▼──────────┐
                                    │      RabbitMQ       │
                                    │  exchange: tasks    │
                                    │  queue: tasks_queue │
                                    └────────────────────┘
```

**Three independent binaries:**

- `cmd/api` — HTTP REST server
- `cmd/worker` — scheduler + consumer + executor
- `cmd/client` — CLI to interact with the API

-----

## Tech Stack

|Component|Technology               |
|---------|-------------------------|
|Language |Go 1.21+                 |
|HTTP     |Gin                      |
|Storage  |Redis (via go-redis/v9)  |
|Messaging|RabbitMQ (via amqp091-go)|
|Docs     |Swagger (swaggo)         |
|Runtime  |Docker + Docker Compose  |

-----

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose

### 1. Start infrastructure

```bash
docker-compose up -d
```

This starts Redis on `:6379` and RabbitMQ on `:5672` (management UI on `:15672`).

### 2. Run the API

```bash
go run ./cmd/api
```

### 3. Run the Worker

```bash
go run ./cmd/worker
```

### 4. Use the CLI

```bash
# Create a task
go run ./cmd/client create --id my-task-1 --name "My Task" --command "echo hello"

# List all tasks
go run ./cmd/client list

# Get a specific task
go run ./cmd/client get --id my-task-1

# Schedule a task
go run ./cmd/client schedule --id my-task-1 --at "2025-12-01T15:00:00Z"

# Execute a task immediately
go run ./cmd/client execute --id my-task-1

# Update a task
go run ./cmd/client update --id my-task-1 --name "Updated Name"

# Delete a task
go run ./cmd/client delete --id my-task-1
```

-----

## Environment Variables

|Variable       |Default                               |Description            |
|---------------|--------------------------------------|-----------------------|
|`REDIS_ADDRESS`|`localhost:6379`                      |Redis server address   |
|`REDIS_MQ_URL` |`amqp://user:password@localhost:5672/`|RabbitMQ connection URL|
|`SERVER_PORT`  |`8080`                                |HTTP API port          |

-----

## API Endpoints

|Method|Path                 |Description                     |
|------|---------------------|--------------------------------|
|GET   |`/health`            |Health check                    |
|POST  |`/tasks`             |Create a task                   |
|GET   |`/tasks`             |List tasks (optional `?status=`)|
|GET   |`/tasks/:id`         |Get task by ID                  |
|PUT   |`/tasks/:id`         |Update task                     |
|DELETE|`/tasks/:id`         |Delete task                     |
|POST  |`/tasks/:id/execute` |Trigger immediate execution     |
|POST  |`/tasks/:id/schedule`|Schedule task for future        |

Swagger docs available at `/swagger/index.html` when running the API.

-----

## Task Lifecycle

```
created → pending → running → completed
                           ↘ failed
```

Tasks with a `scheduled_at` time are stored in a Redis Sorted Set (`scheduled_tasks`) scored by Unix timestamp. The worker polls this set every 30 seconds, publishes due tasks to RabbitMQ, and the consumer picks them up for execution via `sh -c`.

-----

## Project Structure

```
task_scheduler/
├── cmd/
│   ├── api/        # API binary entrypoint
│   ├── client/     # CLI binary + commands
│   └── worker/     # Worker binary entrypoint
├── internal/
│   ├── api/        # HTTP server, routes, handlers
│   ├── domain/     # Task entity + validation
│   ├── messaging/  # RabbitMQ implementation
│   ├── repository/ # TaskHandler interface + Redis impl
│   └── worker/     # Scheduler, consumer, processor
├── pkg/
│   └── config/     # Environment configuration
├── docs/           # Swagger generated docs
└── docker-compose.yml
```

-----

## Running Tests

```bash
go test ./...
```

-----

## Known Limitations & Future Work

- Redis is used as the primary store; data is lost if Redis restarts without persistence configured
- The polling window (30s, 5min lookahead) can cause slight delays in task execution
- `sh -c` execution is not sandboxed; commands run with the worker process’s privileges
- No retry count limit on failed tasks
- No authentication on the HTTP API

See the [Architecture Decisions](#architecture-decisions) section for open design questions.
