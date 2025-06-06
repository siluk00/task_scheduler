package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/siluk00/task_scheduler/internal/domain"
)

const (
	taskKeyPrefix = "task:"
	taskIndex     = "tasks"
)

type TaskRepository struct {
	client *redis.Client
}

// Creates the database connection and returns a new instance of TaskRepository.
// NewTaskRepository initializes a new TaskRepository with the provided Redis client.
func NewTaskRepository(client *redis.Client) *TaskRepository {
	return &TaskRepository{
		client: client,
	}
}

// The TaskRepository struct implements the TaskRepository interface.

// Adds a new task to the Redis database.
func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	// Implementation for creating a task in Redis
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	//Usando pipeline para multplas operações
	//-> Pipeline permite agrupar várias operações Redis em uma única solicitação,
	//o que pode melhorar o desempenho ao reduzir a latência de rede.
	//-> TxPipeline é uma forma de executar comandos Redis em transação,
	//garantindo que todos os comandos sejam executados como uma única unidade atômica.
	pipe := r.client.TxPipeline()
	//Armazena a tarefa no Redis com uma chave única
	// pipe.Set recebe o contexto, a chave, os dados e o tempo de expiração
	// (0 significa sem expiração)
	pipe.Set(ctx, getTaskKey(task.ID), data, 0)
	// Adiciona a tarefa ao índice de tarefas
	// pipe.SAdd recebe o contexto, o nome do índice e o ID da tarefa
	pipe.SAdd(ctx, taskIndex, task.ID)

	if !task.ScheduledAt.IsZero() {
		//pipe.zadd adds the task to the sorted set for scheduled tasks
		// with the score being the scheduled time in Unix format.
		// pipe.ZAdd recebe o contexto, o nome do conjunto ordenado e um objeto Z com o score e o membro
		// O score é o tempo agendado em formato Unix e o membro é o ID da tarefa.
		pipe.ZAdd(ctx, "scheduled_tasks", redis.Z{
			//the Z struct is used to represent a member of a sorted set in Redis.
			Score:  float64(task.ScheduledAt.Unix()),
			Member: task.ID,
		})
	}

	//exec executes all commands in the pipeline atomically.
	_, err = pipe.Exec(ctx)
	return err
}

func (r *TaskRepository) FindById(ctx context.Context, id string) (*domain.Task, error) {
	// Get method retrieves the task data from Redis using the task ID.
	// It returns the task if found, or an error if not found or if there is a Redis error.
	// the Result method returns the value stored at the key.
	// Result is needed because Get returns a *StringCmd,
	// which is a command that can be executed to get the value. The way Redis works...
	data, err := r.client.Get(ctx, getTaskKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // No error, but no task found
		}
		return nil, fmt.Errorf("failed to get task %w from redis", err)
	}

	var task domain.Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task data: %w", err)
	}
	return &task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now()
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}
	// Update the task in Redis	to no-expire
	return r.client.Set(ctx, getTaskKey(task.ID), data, 0).Err()
}

func (r *TaskRepository) List(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error) {
	// This implementation is not enough for big data sets.
	// TODO: considering using SCAN or different data structures for better performance.

	ids, err := r.client.SMembers(ctx, taskIndex).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	var tasks []*domain.Task

	for _, id := range ids {
		task, err := r.FindById(ctx, id)
		if err != nil {
			return nil, err
		}
		if task != nil && (status == "" || task.Status == status) {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	pipe := r.client.TxPipeline()
	pipe.Del(ctx, getTaskKey(id))         // Remove the task by its key
	pipe.SRem(ctx, taskIndex, id)         // Remove the task ID from the index
	pipe.ZRem(ctx, "scheduled_tasks", id) // Remove the task from the scheduled tasks sorted set
	_, err := pipe.Exec(ctx)
	return err
}

func (r *TaskRepository) FindScheduled(ctx context.Context, from, to time.Time) ([]*domain.Task, error) {
	//zrangeByScore retrieves the IDs of tasks scheduled between the specified time range.
	// It uses the ZRangeByScore command to get the members of the sorted set "scheduled_tasks"
	// that have a score (scheduled time) between the Unix timestamps of 'from' and 'to'.
	// The options parameter allows specifying the minimum and maximum scores to filter the results.
	ids, err := r.client.ZRangeByScore(ctx, "scheduled_tasks", &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", from.Unix()),
		Max: fmt.Sprintf("%d", to.Unix()),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to find scheduled tasks: %w", err)
	}

	var tasks []*domain.Task

	for _, id := range ids {
		task, err := r.FindById(ctx, id)
		if err != nil {
			return nil, err
		}
		if task != nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func getTaskKey(id string) string {
	return taskKeyPrefix + id
}
