# Karma Task Scheduler

Scalable task scheduling system using Go, Redis, and RabbitMQ.

## Motivation
I wanted to do myself an application that could use lots of backend concepts in just one software, the best option i found was the task-scheduler since i can use message broker, cache, database knowledge and containers

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

## Quick Start

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
./bin/client create --name "Backup" --command "sudo apt-get update && sudo apt-get upgrade -y"
```

## Usage

## To implement next
- PostgreSQL for persistence
- Architecture Design
- Linear Algebra Tasks 
- Gob faster encoding/decoding
- Two more workers to distribute tasks 
- Use viper for environment variables
- Implement flags for task repetions "daily", "hourly", ...

## Contributing
