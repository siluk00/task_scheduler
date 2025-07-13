# Task Scheduler

Scalable task scheduling system using Go, Redis, and RabbitMQ.

## Key Features
- 🕒 Custom command scheduling
- 📈 Distributed asynchronous execution
- 🔁 Automatic failure retry
- 🔒 JWT authentication
- 📊 Real-time metrics

## Technologies
- **Backend**: Go 1.20+
- **Storage**: Redis
- **Queues**: RabbitMQ
- **Orchestration**: Docker/Kubernetes

## How to Run

### Prerequisites
- Docker and Docker Compose
- Go 1.20+

### Step by Step
```bash
# 1. Clone repository
git clone https://github.com/seu-usuario/task-scheduler.git
cd task-scheduler

# 2. Start infrastructure
docker-compose up -d redis rabbitmq

# 3. Build and run
go build -o bin/api cmd/api/main.go
go build -o bin/worker cmd/worker/main.go

./bin/api
./bin/worker

# 4. Use CLI client
go build -o bin/client cmd/client/main.go
./bin/client create --name "Backup" --command "pg_dump mydb"
```

## To implement next
- PostgreSQL for persistence
- Architecture Design
- Linear Algebra Tasks 
- Gob faster encoding/decoding
- Two more workers to distribute tasks 
