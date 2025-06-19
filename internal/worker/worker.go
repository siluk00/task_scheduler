package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/siluk00/task_scheduler/internal/domain"
	"github.com/siluk00/task_scheduler/internal/messaging/rabbitmq"
	"github.com/siluk00/task_scheduler/internal/repository"
	redisL "github.com/siluk00/task_scheduler/internal/repository/redis"
	"github.com/siluk00/task_scheduler/pkg/config"
)

type TaskWorker struct {
	config      *config.AppConfig
	redisClient *redis.Client
	taskRepo    repository.TaskRepository //CRUD interface of the server
	msgQueue    rabbitmq.MesssageQueue
	running     bool
}

func NewTaskWorker(cfg *config.AppConfig) (*TaskWorker, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: "",
		DB:       0,
	})

	//verifies redis connetion
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err.Error())
	}

	taskRepo := redisL.NewTaskRepository(rdb)

	//gets a connection and a channel in rabbitmq
	msgQueue, err := rabbitmq.NewRabbitMQ(cfg.RedisMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err.Error())
	}

	return &TaskWorker{
		config:      cfg,
		redisClient: rdb,
		taskRepo:    taskRepo,
		msgQueue:    msgQueue,
	}, nil
}

func (w *TaskWorker) Start(ctx context.Context) error {
	w.running = true
	log.Println("Worker started")

	if err := w.SetupRabbitMQ(); err != nil {
		return fmt.Errorf("failed to setup rabbitmq: %w", err)
	}

	go func() {
		if err := w.StartConsumer(ctx); err != nil {
			log.Printf("Consumer stopped with error: %v", err)
			w.Stop(ctx)
		}
	}()

	for w.running {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := w.proccessTasks(ctx); err != nil {
				log.Printf("Error processing tasks: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}

	return nil
}

func (w *TaskWorker) SetupRabbitMQ() error {
	// declare two exchanges, tasks and direct
	err := w.msgQueue.DeclareExchange("tasks", "direct")
	if err != nil {
		return err
	}

	_, err = w.msgQueue.DeclareQueue("tasks_queue")
	if err != nil {
		return nil
	}

	//binds the tasks_queue to the tasks exchange
	return w.msgQueue.BindQueue("tasks_queue", "tasks", "tasks.routing.key")
}

func (w *TaskWorker) proccessTasks(ctx context.Context) error {
	now := time.Now()
	windowEnd := now.Add(5 * time.Minute)

	tasks, err := w.taskRepo.FindScheduled(ctx, now, windowEnd)
	if err != nil {
		return fmt.Errorf("error finding scheduled tasks: %v", err)
	}

	for _, task := range tasks {
		if task.ScheduledAt.After(now) {
			continue
		}

		task.Status = domain.TaskStatusRunning
		if err := w.taskRepo.Update(ctx, task); err != nil {
			log.Printf("Failed to update task %s to running: %v", task.ID, err)
			continue
		}

		if err := w.publishTask(task); err != nil {
			log.Printf("Failed to publish task %s: %v", task.ID, err)
			task.Status = domain.TaskStatusPending
			_ = w.taskRepo.Update(ctx, task)
			continue
		}

		log.Printf("Task %s published for execution", task.ID)
	}

	time.Sleep(30 * time.Second)
	return nil
}

func (w *TaskWorker) publishTask(task *domain.Task) error {
	taskData, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to maarshal task: %v", err)
	}

	return w.msgQueue.Publish("tasks", "task.routing.key", taskData)
}

func (w *TaskWorker) Stop(ctx context.Context) {
	w.running = false
	if err := w.msgQueue.CLose(); err != nil {
		log.Printf("Error closing message queue: %v", err)
	}
}
